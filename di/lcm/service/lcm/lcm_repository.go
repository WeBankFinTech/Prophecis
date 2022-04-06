/*
 * Copyright 2017-2018 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lcm

import (
	"errors"
	"strconv"
	"webank/DI/commons/constants"
	"webank/DI/commons/logger"
	"webank/DI/trainer/trainer/grpc_trainer_v2"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"strings"
	"webank/DI/trainer/trainer"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	TRAINING_JOB_QUEUE_NS = "TRAINING_JOB_QUEUE_NS"
	TRAINING_JOB_QUEUE_   = "TRAINING_JOB_QUEUE_"
)

type trainingsRepository struct {
	session    *mgo.Session
	database   string
	collection string
}

type repository interface {
	Store(c *trainer.TrainingRecord) error
	Find(trainingID string) (*trainer.TrainingRecord, error)
	FindTrainingStatus(trainingID string) (*grpc_trainer_v2.TrainingStatus, error)
	FindTrainingStatusID(trainingID string) (grpc_trainer_v2.Status, error)
	FindTrainingSummaryMetricsString(trainingID string) (string, error)
	// FIXME MLSS Change: get models filter by username and namespace
	FindAll(userID string, page string, size string) (*trainer.TrainingRecordPaged, error)
	FindAllByUserIdAndNamespace(user *string, namespace *string, page string, size string) (*trainer.TrainingRecordPaged, error)
	FindAllByUserIdAndNamespaceList(user *string, namespace *[]string, page string, size string) (*trainer.TrainingRecordPaged, error)
	FindCurrentlyRunningTrainings(limit int) ([]*trainer.TrainingRecord, error)
	Delete(trainingID string) error
	Close()
}

type jobHistoryRepository interface {
	RecordJobStatus(e *trainer.JobHistoryEntry) error
	GetJobStatusHistory(trainingID string) []*trainer.JobHistoryEntry
	Close()
}

// newTrainingsRepository creates a new training repo for storing training data. The mongo URI can contain all the necessary
// connection information. See here: http://docs.mongodb.org/manual/reference/connection-string/
// However, we also support not putting the username/password in the connection URL and provide is separately.
func newTrainingsRepository(mongoURI string, database string, username string, password string, authenticationDatabase string,
	cert string, collection string) (repository, error) {
	log := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	log.Debugf("Creating mongo training repository for %s, collection %s:", mongoURI, collection)

	session, err := trainer.ConnectMongo(mongoURI, database, username, password, authenticationDatabase, cert)
	if err != nil {
		return nil, err
	}
	collectionObj := session.DB(database).C(collection)

	repo := &trainingsRepository{
		session:    session,
		database:   collectionObj.Database.Name,
		collection: collection,
	}

	// create index
	collectionObj.EnsureIndexKey("user_id", "training_id")

	return repo, nil
}

func NewTrainingsRepository(mongoURI string, database string, username string, password string, authenticationDatabase string,
	cert string, collection string) (repository, error) {
	log := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	log.Debugf("Creating mongo training repository for %s, collection %s:", mongoURI, collection)

	session, err := trainer.ConnectMongo(mongoURI, database, username, password, authenticationDatabase, cert)
	if err != nil {
		return nil, err
	}
	collectionObj := session.DB(database).C(collection)

	repo := &trainingsRepository{
		session:    session,
		database:   collectionObj.Database.Name,
		collection: collection,
	}

	// create index
	collectionObj.EnsureIndexKey("user_id", "training_id")

	return repo, nil
}

// FIXME MLSS Change: v_1.4.1 add func to get queues from mongoDB
func getQueuesFromDB(mongoURI string, database string, username string, password string, authenticationDatabase string,
	cert string) (q map[string]*queueHandler, err error) {
	log := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	log.Debugf("get queues from mongo")

	session, err := trainer.ConnectMongo(mongoURI, database, username, password, authenticationDatabase, cert)
	if err != nil {
		return nil, err
	}

	queues := make(map[string]*queueHandler)
	collectionNames, err := session.DB(database).CollectionNames()
	for _, collectionName := range collectionNames {
		if strings.HasPrefix(collectionName, TRAINING_JOB_QUEUE_NS) {
			namespaceName := strings.Replace(collectionName, TRAINING_JOB_QUEUE_, "", 1)
			queue, err := trainer.GetNewTrainingJobQueue(mongoURI, database, username, password, authenticationDatabase,
				cert, trainer.QueueName(namespaceName), trainer.LockName(namespaceName))
			if err != nil {
				log.WithError(err).Fatalf("Cannot create queue with %s %s %s", viper.GetString(MongoAddressKey),
					viper.GetString(MongoDatabaseKey), viper.GetString(MongoUsernameKey))
				return nil, err
			}
			queues[collectionName] = &queueHandler{make(chan struct{}), queue, constants.Fault}
		}
	}
	return queues, nil
}

// newJobHistoryRepository creates a new repo for storing job status history entries.
func newJobHistoryRepository(mongoURI string, database string, username string, password string, authenticationDatabase string,
	cert string, collection string) (jobHistoryRepository, error) {
	log := logger.LocLogger(log.StandardLogger().WithField("module", "jobHistoryRepository"))
	log.Debugf("Creating mongo repository for %s, collection %s:", mongoURI, collection)

	session, err := trainer.ConnectMongo(mongoURI, database, username, password, authenticationDatabase,cert)
	if err != nil {
		return nil, err
	}
	collectionObj := session.DB(database).C(collection)

	repo := &trainingsRepository{
		session:    session,
		database:   collectionObj.Database.Name,
		collection: collection,
	}
	return repo, err
}

func (r *trainingsRepository) Store(t *trainer.TrainingRecord) error {
	sess := r.session.Clone()
	defer sess.Close()

	var err error
	if t.ID == "" {
		logger.GetLogger().Info("insert training record：", t)
		err = sess.DB(r.database).C(r.collection).Insert(t)
	} else {
		logger.GetLogger().Info("update training record：", t)
		err = sess.DB(r.database).C(r.collection).Update(bson.M{"_id": t.ID}, t)
	}
	if err != nil {
		logWith(t.TrainingID, t.UserID).Errorf("Error storing training record: %s", err.Error())
		return err
	}

	return nil
}

func (r *trainingsRepository) Find(trainingID string) (*trainer.TrainingRecord, error) {
	tr := &trainer.TrainingRecord{}
	sess := r.session.Clone()
	defer sess.Close()
	err := r.queryDatabase(&bson.M{"training_id": trainingID}, sess).One(tr)
	if err != nil {
		logWithTraining(trainingID).WithError(err).Debugf("Cannot retrieve training record")
		return nil, err
	}
	return tr, nil
}

func (r *trainingsRepository) FindTrainingStatus(trainingID string) (*grpc_trainer_v2.TrainingStatus, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository").WithField(logger.LogkeyTrainingID, trainingID))
	sess := r.session.Clone()
	defer sess.Close()

	tr := &grpc_trainer_v2.TrainingStatus{}
	err := r.queryDatabase(&bson.M{"training_id": trainingID}, sess).Select(bson.M{"TrainingStatus": 1}).One(tr)
	if err != nil {
		logr.WithError(err).Debugf("Cannot retrieve training record")
		return nil, err
	}

	return tr, nil
}

func (r *trainingsRepository) FindTrainingStatusID(trainingID string) (grpc_trainer_v2.Status, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository").WithField(logger.LogkeyTrainingID, trainingID))
	sess := r.session.Clone()
	defer sess.Close()

	tr := &trainer.TrainingRecord{}
	err := r.queryDatabase(&bson.M{"training_id": trainingID}, sess).One(tr)
	if err != nil {
		logWithTraining(trainingID).WithError(err).Debugf("Cannot retrieve training record")
		return -1, err
	}

	if tr.TrainingStatus != nil {
		return tr.TrainingStatus.Status, nil
	}
	logr.Debugf("Status not found")
	return grpc_trainer_v2.Status_NOT_STARTED, nil
}

func (r *trainingsRepository) FindTrainingSummaryMetricsString(trainingID string) (string, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository").WithField(logger.LogkeyTrainingID, trainingID))
	sess := r.session.Clone()
	defer sess.Close()

	tr := &trainer.TrainingRecord{}
	err := r.queryDatabase(&bson.M{"training_id": trainingID}, sess).One(tr)
	if err != nil {
		logWithTraining(trainingID).WithError(err).Debugf("Cannot retrieve training record")
		return "", err
	}

	if tr.TrainingStatus != nil {
		return tr.Metrics.String(), nil
	}
	logr.Debugf("Status not found")
	return "", nil
}

func (r *trainingsRepository) FindAll(userID string, page string, size string) (*trainer.TrainingRecordPaged, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	var tr []*trainer.TrainingRecord
	sess := r.session.Clone()
	defer sess.Close()

	//query statement
	query := r.queryDatabase(&bson.M{"user_id": userID}, sess).Sort("-training_status.submissiontimestamp")
	countQuery := r.queryDatabase(&bson.M{"user_id": userID}, sess).Sort("-training_status.submissiontimestamp")

	//paged
	var sizeInt int
	if page != "" && size != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return nil, errors.New(" page parse to int failed.")
		}
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			return nil, errors.New(" size parse to int failed.")
		}
		offset := (pageInt - 1) * sizeInt
		query = query.Skip(offset).Limit(sizeInt)
	}
	err := query.All(&tr)
	if err != nil {
		log.WithField(logger.LogkeyUserID, userID).Errorf("Cannot retrieve all training records: %s", err.Error())
		return nil, err
	}

	//count jobs num
	total, err := countQuery.Count()
	if err != nil {
		log.WithField(logger.LogkeyUserID, userID).Errorf("Cannot count all training records: %s", err.Error())
		return nil, err
	}
	//get pages
	var pages int
	if sizeInt != 0 {
		pages = total / sizeInt
		if total%sizeInt != 0 {
			pages += 1
		}
	}
	logr.Infof("total: %v, pages: %v", total, pages)
	paged := trainer.TrainingRecordPaged{
		TrainingRecords: tr,
		Pages:           pages,
		Total:           total,
	}

	return &paged, nil
}

// FIXME MLSS Change: get models filter by username and namespace
func (r *trainingsRepository) FindAllByUserIdAndNamespace(user *string, namespace *string, page string, size string) (*trainer.TrainingRecordPaged, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	var tr []*trainer.TrainingRecord
	sess := r.session.Clone()
	defer sess.Close()

	m := &bson.M{}
	if *user != "" {
		(*m)["user_id"] = user
	}
	if *namespace != "" {
		(*m)["namespace"] = namespace
	}
	// FIXME MLSS Change: get models filter by username and namespace
	query := r.queryDatabase(m, sess).Sort("-training_status.submissiontimestamp")
	countQuery := r.queryDatabase(m, sess).Sort("-training_status.submissiontimestamp")

	//paged
	var sizeInt int
	logr.Infof("var sizeInt int: %v", sizeInt)
	if page != "" && size != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return nil, errors.New(" page parse to int failed.")
		}
		sizeInt, err = strconv.Atoi(size)
		if err != nil {
			return nil, errors.New(" size parse to int failed.")
		}
		offset := (pageInt - 1) * sizeInt
		query = query.Skip(offset).Limit(sizeInt)
	}
	logr.Infof("strconv.Atoi(size): %v", sizeInt)

	err := query.All(&tr)
	if err != nil {
		log.WithField(logger.LogkeyUserID, user).Errorf("Cannot retrieve all training records: %s", err.Error())
		return nil, err
	}

	//count jobs num
	total, err := countQuery.Count()
	if err != nil {
		log.WithField(logger.LogkeyUserID, user).Errorf("Cannot count all training records: %s", err.Error())
		return nil, err
	}
	//get pages
	var pages int
	if sizeInt != 0 {
		pages = total / sizeInt
		logr.Infof("pages = total / sizeInt: %v", pages)

		if total%sizeInt != 0 {
			pages += 1
		}
	}
	logr.Infof("total: %v, pages: %v", total, pages)
	paged := trainer.TrainingRecordPaged{
		TrainingRecords: tr,
		Pages:           pages,
		Total:           total,
	}

	return &paged, nil
}

// FIXME MLSS Change: get models filter by username and namespace list
func (r *trainingsRepository) FindAllByUserIdAndNamespaceList(user *string, namespaces *[]string, page string, size string) (*trainer.TrainingRecordPaged, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	var tr []*trainer.TrainingRecord
	sess := r.session.Clone()
	defer sess.Close()
	//create query object
	m := &bson.M{}
	var n []bson.M
	if *user != "" {
		//(*m)["user_id"] = *user
		n = append(n, bson.M{"user_id": *user})
	}
	if len(*namespaces) > 0 {
		//var n []bson.M
		for _, v := range *namespaces {
			n = append(n, bson.M{"namespace": v})
		}
	}
	logr.Infof("FindAllByUserIdAndNamespaceList query info: %v", *m)
	if len(n) != 0 {
		(*m)["$or"] = n
	}

	query := r.queryDatabase(m, sess).Sort("-training_status.submissiontimestamp")
	countQuery := r.queryDatabase(m, sess).Sort("-training_status.submissiontimestamp")
	//paged
	var sizeInt int
	if page != "" && size != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return nil, errors.New(" page parse to int failed.")
		}
		sizeInt, err = strconv.Atoi(size)
		if err != nil {
			return nil, errors.New(" size parse to int failed.")
		}
		offset := (pageInt - 1) * sizeInt
		query = query.Skip(offset).Limit(sizeInt)
	}
	//err = query.Distinct("_id", &tr)
	err := query.All(&tr)
	if err != nil {
		log.WithField(logger.LogkeyUserID, user).Errorf("Cannot retrieve all training records: %s", err.Error())
		return nil, err
	}
	//count jobs num
	total, err := countQuery.Count()
	if err != nil {
		log.WithField(logger.LogkeyUserID, user).Errorf("Cannot count all training records: %s", err.Error())
		return nil, err
	}
	//get pages
	var pages int
	if sizeInt != 0 {
		pages = total / sizeInt
		if total%sizeInt != 0 {
			pages += 1
		}
	}
	logr.Infof("total: %v, pages: %v, tr: s%", total, pages, tr)
	paged := trainer.TrainingRecordPaged{
		TrainingRecords: tr,
		Pages:           pages,
		Total:           total,
	}
	return &paged, nil
}

func (r *trainingsRepository) Delete(trainingID string) error {
	sess := r.session.Clone()
	defer sess.Close()
	// Perform a soft delete: retain only non-sensitive details of the training record. Note: instead of
	// deleting fields from the record, we upsert with a new record to specify explicitly what should be retained.

	// 1. fetch the record
	existing, err := r.Find(trainingID)
	if err != nil {
		logWithTraining(trainingID).WithError(err).Debugf("Unable to find training record for (soft-)deletion, ID %s: %s", trainingID, err)
		return err
	}
	// 2. update the record
	selector := bson.M{"training_id": trainingID}
	var resources *grpc_trainer_v2.ResourceRequirements
	var status grpc_trainer_v2.TrainingStatus
	var framework *grpc_trainer_v2.Framework
	if existing.Training != nil {
		resources = existing.Training.Resources
	}
	if existing.TrainingStatus != nil {
		status = *existing.TrainingStatus
	}
	if existing.ModelDefinition != nil {
		framework = existing.ModelDefinition.Framework
	}
	newRecord := &trainer.TrainingRecord{
		TrainingID: trainingID,
		UserID:     existing.UserID,
		JobID:      existing.JobID,
		ModelDefinition: &grpc_trainer_v2.ModelDefinition{
			Framework: framework,
		},
		Training: &grpc_trainer_v2.Training{
			Resources: resources,
		},
		TrainingStatus: &grpc_trainer_v2.TrainingStatus{
			Status:                 status.Status,
			ErrorCode:              status.ErrorCode,
			SubmissionTimestamp:    status.SubmissionTimestamp,
			CompletionTimestamp:    status.CompletionTimestamp,
			DownloadStartTimestamp: status.DownloadStartTimestamp,
			ProcessStartTimestamp:  status.ProcessStartTimestamp,
			StoreStartTimestamp:    status.StoreStartTimestamp,
		},
		Deleted: true,
	}
	_, err1 := sess.DB(r.database).C(r.collection).Upsert(selector, newRecord)
	if err1 != nil {
		logWithTraining(trainingID).Errorf("Cannot (soft-)delete training record: %s", err.Error())
		return err1
	}
	return nil
}

// queryDatabase serves as the single entry point to run DB queries for this repository. It takes as parameter
// a selector to use for MongoDB's Find(...) method, and returns the query result. Importantly, the method appends
// a "deleted" flag to the query selector to make sure we are never returning records that have been soft-deleted.
func (r *trainingsRepository) queryDatabase(selector *bson.M, sess *mgo.Session) *mgo.Query {
	if selector == nil {
		selector = &bson.M{}
	}
	if (*selector)["deleted"] == nil {
		(*selector)["deleted"] = bson.M{"$ne": true}
	}
	return sess.DB(r.database).C(r.collection).Find(selector)
}

func (r *trainingsRepository) FindCurrentlyRunningTrainings(limit int) ([]*trainer.TrainingRecord, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	sess := r.session.Clone()
	defer sess.Close()

	logr.Debugf("find currently running trainings using database: %s collection: %s live servers: %s", r.database, r.collection, sess.LiveServers())

	var tr []*trainer.TrainingRecord
	//sorting by id in descending fashion(hence the - before id), assumption being records in mongo are being created with auto generated id which has a notion of timestamp built into it
	err := r.queryDatabase(nil, sess).Sort("-_id").Limit(limit).Select(bson.M{"training_status": 1, "training.resources": 2, "training_id": 3}).All(&tr)
	return tr, err
}

func (r *trainingsRepository) RecordJobStatus(e *trainer.JobHistoryEntry) error {
	sess := r.session.Clone()
	defer sess.Close()

	setOnInsert := make(map[string]interface{})
	setOnInsert["$setOnInsert"] = e
	_, err := sess.DB(r.database).C(r.collection).Upsert(e, setOnInsert)
	if err != nil {
		logWithTraining(e.TrainingID).Errorf("Error storing job history entry: %s", err.Error())
		return err
	}

	return nil
}

func (r *trainingsRepository) GetJobStatusHistory(trainingID string) []*trainer.JobHistoryEntry {
	sess := r.session.Clone()
	defer sess.Close()
	selector := bson.M{}
	if trainingID != "" {
		selector["training_id"] = trainingID
	}
	var result []*trainer.JobHistoryEntry
	sess.DB(r.database).C(r.collection).Find(selector).Sort("timestamp").All(&result)
	return result
}

func (r *trainingsRepository) Close() {
	log.Debugf("Closing mongo session")
	defer r.session.Close()
}
