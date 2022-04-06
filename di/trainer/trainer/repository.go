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

package trainer

import (
	"errors"
	"strconv"
	"webank/DI/commons/constants"
	"webank/DI/commons/logger"
	"webank/DI/trainer/trainer/grpc_trainer_v2"

	log "github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// TrainingRecord is the data structure we store in the Mongo collection "training_jobs"
type TrainingRecord struct {
	ID                    bson.ObjectId                    `bson:"_id,omitempty" json:"id"`
	TrainingID            string                           `bson:"training_id" json:"training_id"`
	UserID                string                           `bson:"user_id" json:"user_id"`
	JobID                 string                           `bson:"job_id" json:"job_id"`
	ModelDefinition       *grpc_trainer_v2.ModelDefinition `bson:"model_definition,omitempty" json:"model_definition"`
	Training              *grpc_trainer_v2.Training        `bson:"training,omitempty" json:"training"`
	Datastores            []*grpc_trainer_v2.Datastore     `bson:"data_stores,omitempty" json:"data_stores"`
	TrainingStatus        *grpc_trainer_v2.TrainingStatus  `bson:"training_status,omitempty" json:"training_status"`
	Metrics               *grpc_trainer_v2.Metrics         `bson:"metrics,omitempty" json:"metrics"`
	Deleted               bool                             `bson:"deleted,omitempty" json:"deleted"`
	EvaluationMetricsSpec string                           `bson:"evaluation_metrics_spec,omitempty" json:"evaluation_metrics_spec"`
	Namespace             string                           `bson:"namespace,omitempty" json:"namespace"`
	GID                   string                           `bson:"gid" json:"gid"`
	UID                   string                           `bson:"uid" json:"uid"`
	GuardianToken         string                           `bson:"guardian_token" json:"guardian_token"`
	JobAlert              string                           `bson:"job_alert,omitempty" json:"job_alert,omitempty"`
	PSs                   string                           `bson:"pss,omitempty" json:"pss,omitempty"`
	PSCPU                 string                           `bson:"ps_cpu,omitempty" json:"ps_cpu,omitempty"`
	PSImage               string                           `bson:"ps_image,omitempty" json:"ps_image,omitempty"`
	PSMemory              string                           `bson:"ps_memory,omitempty" json:"ps_memory,omitempty"`
	CodeSelector          string                           `bson:"code_selector,omitempty" json:"code_selector"`
	DataPath              string                           `bson:"data_path,omitempty" json:"data_path"`

	StoreStatus string `bson:"store_status,omitempty" json:"store_status"`
	JobType     string `bson:"job_type,omitempty" json:"job_type"`
	TFosRequest string `bson:"tfos_request,omitempty" json:"tfos_request"`

	LinkisExecId string `bson:"linkis_exec_id,omitempty" json:"linkis_exec_id"`
	ExpRunId     string `bson:"exp_run_id,omitempty" json:"exp_run_id"`
	ExpName      string `bson:"exp_name,omitempty" json:"exp_name"`
	FileName     string `bson:"file_name,omitempty" json:"file_name"`
	FilePath     string `bson:"file_path,omitempty" json:"file_path"`

	ProxyUser string `bson:"proxy_user,omitempty" json:"proxy_user"`

	DataSet   *grpc_trainer_v2.DataSet `bson:"data_set,omitempty" json:"data_set"`
	JobParams string `bson:"job_params,omitempty" json:"job_params"`
	MFModel   *grpc_trainer_v2.MFModel `bson:"mf_model,omitempty" json:"mf_model"`
	Algorithm  string `bson:"algorithm,omitempty" json:"algorithm"`
	FitParams  string `bson:"fit_params,omitempty" json:"fit_params"`
	APIType   string  `bson:"API_type,omitempty" json:"API_type"`
}

type DataSet struct {
	TrainingDataPath    string `bson:"job_params,omitempty" json:"training_data_path"`
	TrainingLabelPath   string `bson:"training_label_path,omitempty" json:"training_label_path,omitempty"`
	TestingDataPath     string `bson:"testing_data_path,omitempty" json:"testing_data_path,omitempty"`
	TestingLabelPath    string `bson:"testing_label_path,omitempty" json:"testing_label_path,omitempty"`
	ValidationDataPath  string `bson:"validation_data_path,omitempty" json:"validation_data_path,omitempty"`
	ValidationLabelPath string `bson:"validation_label_path,omitempty" json:"validation_label_path,omitempty"`
}

// FIXME MLSS Change: get models with page info
type TrainingRecordPaged struct {
	TrainingRecords []*TrainingRecord `bson:"training_records" json:"training_records"`
	Pages           int               `bson:"pages" json:"pages"`
	Total           int               `bson:"total" json:"total"`
}

// JobHistoryEntry stores training job status history in the Mongo collection "job_history"
type JobHistoryEntry struct {
	ID            bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	TrainingID    string                 `bson:"training_id" json:"training_id"`
	Timestamp     string                 `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	Status        grpc_trainer_v2.Status `bson:"status,omitempty" json:"status,omitempty"`
	StatusMessage string                 `bson:"status_message,omitempty" json:"status_message,omitempty"`
	ErrorCode     string                 `bson:"error_code,omitempty" json:"error_code,omitempty"`
}

type trainingsRepository struct {
	session    *mgo.Session
	database   string
	collection string
}

type Repository interface {
	Store(c *TrainingRecord) error
	Find(trainingID string) (*TrainingRecord, error)
	FindTrainingStatus(trainingID string) (*grpc_trainer_v2.TrainingStatus, error)
	FindTrainingStatusID(trainingID string) (grpc_trainer_v2.Status, error)
	FindTrainingSummaryMetricsString(trainingID string) (string, error)
	// FIXME MLSS Change: get models filter by username and namespace
	FindAll(userID string, page string, size string) (*TrainingRecordPaged, error)
	FindAllByStatus(status grpc_trainer_v2.Status, page string, size string) (*TrainingRecordPaged, error)
	FindAllByUserIdAndNamespace(user *string, namespace *string, page string, size string) (*TrainingRecordPaged, error)
	FindAllByUserIdAndNamespaceList(user *string, namespace *[]string, page string, size string, clusterName string) (*TrainingRecordPaged, error)
	FindCurrentlyRunningTrainings(limit int) ([]*TrainingRecord, error)
	FindAllNotCompletedLinkisJobs(jobDone bool, logStored bool) ([]*TrainingRecord, error)
	FindStoreStatus(trainingID string) (string, error)
	Delete(trainingID string) error
	Close()
	GetSession() *mgo.Session
}

type jobHistoryRepository interface {
	RecordJobStatus(e *JobHistoryEntry) error
	GetJobStatusHistory(trainingID string) []*JobHistoryEntry
	Close()
}

// newTrainingsRepository creates a new training repo for storing training data. The mongo URI can contain all the necessary
// connection information. See here: http://docs.mongodb.org/manual/reference/connection-string/
// However, we also support not putting the username/password in the connection URL and provide is separately.
func newTrainingsRepository(mongoURI string, database string, username string, password string, authenticationDatabase string,
	cert string, collection string) (Repository, error) {
	log := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	log.Debugf("Creating mongo training Repository for %s, collection %s:", mongoURI,authenticationDatabase, collection)

	session, err := ConnectMongo(mongoURI, database, username, password, authenticationDatabase, cert)
	if err != nil {
		return nil, err
	}
	collectionObj := session.DB(database).C(collection)
	LogLockCollectionObj := session.DB(database).C("LogLock")
	LogLockIndex := mgo.Index{
		Key:    []string{"Node", "TrainingID"},
		Unique: true,
	}
	err = LogLockCollectionObj.EnsureIndex(LogLockIndex)

	repo := &trainingsRepository{
		session:    session,
		database:   collectionObj.Database.Name,
		collection: collection,
	}

	// create index
	err = collectionObj.EnsureIndexKey("user_id", "training_id")

	return repo, nil
}

func NewTrainingsRepository(mongoURI string, database string, username string, password string, authenticationDatabase string,
	cert string, collection string) (Repository, error) {
	log := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	log.Debugf("Creating mongo training Repository for %s, collection %s:", mongoURI, collection)

	session, err := ConnectMongo(mongoURI, database, username, password, authenticationDatabase, cert)
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

// newJobHistoryRepository creates a new repo for storing job status history entries.
func newJobHistoryRepository(mongoURI string, database string, username string, password string,
	authenticationDatabase string, cert string, collection string) (jobHistoryRepository, error) {
	log := logger.LocLogger(log.StandardLogger().WithField("module", "jobHistoryRepository"))
	log.Debugf("Creating mongo Repository for %s, collection %s:", mongoURI, collection)

	session, err := ConnectMongo(mongoURI, database, username, password, authenticationDatabase, cert)
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

func (r *trainingsRepository) Store(t *TrainingRecord) error {
	sess := r.session.Clone()
	defer sess.Close()

	var err error
	if t.ID == "" {
		err = sess.DB(r.database).C(r.collection).Insert(t)
	} else {
		err = sess.DB(r.database).C(r.collection).Update(bson.M{"_id": t.ID}, t)
	}
	if err != nil {
		logWith(t.TrainingID, t.UserID).Errorf("Error storing training record: %s", err.Error())
		return err
	}

	return nil
}


func (r *trainingsRepository) FindAllByStatus(status grpc_trainer_v2.Status, page string, size string) (*TrainingRecordPaged, error){
	sess := r.session.Clone()
	defer sess.Close()
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	var tr []*TrainingRecord

	//query statement
	query := r.queryDatabase(&bson.M{"training_status.status": status}, sess)//.Select()
	countQuery := r.queryDatabase(&bson.M{"training_status.status": status}, sess)

	//err := r.queryDatabase(nil, sess).Sort("-_id").Limit(limit).Select(bson.M{"training_status": 1, "training.resources": 2, "training_id": 3}).All(&tr)


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
		log.Errorf("Cannot retrieve all training records: %s", err.Error())
		return nil, err
	}

	//count jobs num
	total, err := countQuery.Count()
	if err != nil {
		log.Errorf("Cannot count all training records: %s", err.Error())
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
	logr.Debugf("total: %v, pages: %v", total, pages)
	paged := TrainingRecordPaged{
		TrainingRecords: tr,
		Pages:           pages,
		Total:           total,
	}

	return &paged, nil
}


func (r *trainingsRepository) Find(trainingID string) (*TrainingRecord, error) {
	tr := &TrainingRecord{}
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

	tr := &TrainingRecord{}
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

	tr := &TrainingRecord{}
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

func (r *trainingsRepository) FindAll(userID string, page string, size string) (*TrainingRecordPaged, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	var tr []*TrainingRecord
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
	logr.Debugf("total: %v, pages: %v", total, pages)
	paged := TrainingRecordPaged{
		TrainingRecords: tr,
		Pages:           pages,
		Total:           total,
	}

	return &paged, nil
}

// FIXME MLSS Change: get models filter by username and namespace
func (r *trainingsRepository) FindAllByUserIdAndNamespace(user *string, namespace *string, page string, size string) (*TrainingRecordPaged, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	var tr []*TrainingRecord
	sess := r.session.Clone()
	defer sess.Close()

	queryMap := &bson.M{}
	if *user != "" {
		(*queryMap)["user_id"] = user
	}
	if *namespace != "" {
		(*queryMap)["namespace"] = namespace
	}
	// FIXME MLSS Change: get models filter by username and namespace
	query := r.queryDatabase(queryMap, sess).Sort("-training_status.submissiontimestamp")
	countQuery := r.queryDatabase(queryMap, sess).Sort("-training_status.submissiontimestamp")

	//paged
	var sizeInt int
	logr.Debugf("var sizeInt int: %v", sizeInt)
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
	logr.Debugf("strconv.Atoi(size): %v", sizeInt)

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
		logr.Debugf("pages = total / sizeInt: %v", pages)

		if total%sizeInt != 0 {
			pages += 1
		}
	}
	logr.Debugf("total: %v, pages: %v", total, pages)
	paged := TrainingRecordPaged{
		TrainingRecords: tr,
		Pages:           pages,
		Total:           total,
	}

	return &paged, nil
}

// FIXME MLSS Change: get models filter by username and namespace list
func (r *trainingsRepository) FindAllByUserIdAndNamespaceList(user *string, namespaces *[]string, page string, size string, clusterName string) (*TrainingRecordPaged, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	var tr []*TrainingRecord
	sess := r.session.Clone()
	defer sess.Close()
	//create query object
	queryMap := &bson.M{}
	var n []bson.M
	logr.Debugf("FindAllByUserIdAndNamespaceList query info with user: %v, namespaces: %v, clusterName: %v", user, namespaces, clusterName)
	if *user != "" {
		// FIXME MLSS Change: v_1.5.0 add logic to get NB from k8s
		if clusterName == bdap {
			n = append(n, bson.M{"user_id": *user, "namespace": bson.M{"$regex": "^(ns-)([a-z]*-){1}bdap-.*(?<!-safe)$"}})
		}
		if clusterName == bdapsafe {
			n = append(n, bson.M{"user_id": *user, "namespace": bson.M{"$regex": "^(ns-)([a-z]*-){1}bdap-.*(-safe)$"}})
		}
		if clusterName == bdp {
			n = append(n, bson.M{"user_id": *user, "namespace": bson.M{"$regex": "^(ns-)([a-z]*-){1}bdp-.*"}})
		}
		if clusterName == defaultValue {
			n = append(n, bson.M{"user_id": *user})
		}
		//(*queryMap)["user_id"] = *user
	}
	if len(*namespaces) > 0 {
		//var n []bson.M
		for _, v := range *namespaces {
			n = append(n, bson.M{"namespace": v})
		}
	}
	logr.Debugf("FindAllByUserIdAndNamespaceList query info: %v", *queryMap)
	if len(n) != 0 {
		(*queryMap)["$or"] = n
	}
	logr.Debugf("FindAllByUserIdAndNamespaceList query info by n: %v", *queryMap)
	query := r.queryDatabase(queryMap, sess).Sort("-training_status.submissiontimestamp")
	countQuery := r.queryDatabase(queryMap, sess).Sort("-training_status.submissiontimestamp")
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
	logr.Debugf("total: %v, pages: %v, tr: %v", total, pages, tr)
	paged := TrainingRecordPaged{
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
	newRecord := &TrainingRecord{
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

// queryDatabase serves as the single entry point to run DB queries for this Repository. It takes as parameter
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

func (r *trainingsRepository) FindCurrentlyRunningTrainings(limit int) ([]*TrainingRecord, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	sess := r.session.Clone()
	defer sess.Close()

	logr.Debugf("find currently running trainings using database: %s collection: %s live servers: %s", r.database, r.collection, sess.LiveServers())

	var tr []*TrainingRecord
	//sorting by id in descending fashion(hence the - before id), assumption being records in mongo are being created with auto generated id which has a notion of timestamp built into it
	err := r.queryDatabase(nil, sess).Sort("-_id").Limit(limit).Select(bson.M{"training_status": 1, "training.resources": 2, "training_id": 3}).All(&tr)
	return tr, err
}

func (r *trainingsRepository) FindAllNotCompletedLinkisJobs(jobDone, logStored bool) ([]*TrainingRecord, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))
	sess := r.session.Clone()
	defer sess.Close()

	logr.Debugf("find currently running trainings using database: %s collection: %s live servers: %s", r.database, r.collection, sess.LiveServers())

	var tr []*TrainingRecord
	//sorting by id in descending fashion(hence the - before id), assumption being records in mongo are being created with auto generated id which has a notion of timestamp built into it

	ms := bson.M{
		"job_type": constants.TFOS,
	}
	if !logStored {
		ms[constants.LOG_STORING_KEY] = bson.M{"$ne": constants.LOG_STORING_COMPLETED}
	}
	if jobDone {
		ms["$and"] = []bson.M{
			{"training_status.status": bson.M{"$ne": 10}},
			{"training_status.status": bson.M{"$ne": 60}},
		}
	}

	err := r.queryDatabase(&ms, sess).Sort("-_id").All(&tr)
	//n, err := r.queryDatabase(nil, sess).Sort("-_id").Select(bson.M{"job_type": "tfos", "training_status.status": bson.M{"$ne": 10}}).Count()
	//logr.Infof("queryDatabase: %v", n)
	//n2, err := r.queryDatabase(nil, sess).Sort("-_id").Select(bson.M{"job_type": "tfos", "training_status.status": bson.M{"$ne": 60}}).Count()
	//logr.Infof("queryDatabase: %v", n2)
	return tr, err
}

func (r *trainingsRepository) RecordJobStatus(e *JobHistoryEntry) error {
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

func (r *trainingsRepository) GetJobStatusHistory(trainingID string) []*JobHistoryEntry {
	sess := r.session.Clone()
	defer sess.Close()
	selector := bson.M{}
	if trainingID != "" {
		selector["training_id"] = trainingID
	}
	var result []*JobHistoryEntry
	sess.DB(r.database).C(r.collection).Find(selector).Sort("timestamp").All(&result)
	return result
}

func (r *trainingsRepository) Close() {
	log.Debugf("Closing mongo session")
	defer r.session.Close()
}

func (r *trainingsRepository) FindStoreStatus(trainingID string) (string, error) {
	logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository").WithField(logger.LogkeyTrainingID, trainingID))
	sess := r.session.Clone()
	defer sess.Close()

	var status string
	err := r.queryDatabase(&bson.M{"training_id": trainingID}, sess).Select(bson.M{"StoreStatus": 1}).One(status)
	if err != nil {
		logr.WithError(err).Debugf("Cannot retrieve training record")
		return "", err
	}

	return status, nil
}

func (r *trainingsRepository) GetSession() *mgo.Session {
	return r.session
}
