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
	"fmt"
	"testing"
	"time"

	"webank/DI/commons/config"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"webank/DI/trainer/trainer/grpc_trainer_v2"
)

var mongoAddress, mongoDatabase, mongoUsername, mongoPassword, mongoCertLocation string

func init() {
	config.InitViper()

	// mongo default settings
	viper.SetDefault(mongoAddressKey, "mongodb://localhost:27017/dlaas_trainer")
	mongoAddress = viper.GetString(mongoAddressKey)   // overwrite this by setting ENV var DLAAS_MONGO_ADDRESS
	mongoDatabase = viper.GetString(mongoDatabaseKey) // overwrite this by setting ENV var DLAAS_MONGO_ADDRESS
	mongoUsername = viper.GetString(mongoUsernameKey) // overwrite this by setting ENV var DLAAS_MONGO_USERNAME
	mongoPassword = viper.GetString(mongoPasswordKey) // overwrite this by setting ENV var DLAAS_MONGO_PASSWORD
	mongoCertLocation = config.GetMongoCertLocation()
}

func createTrainings() []*TrainingRecord {
	arr := make([]*TrainingRecord, 10)

	for i := 0; i < 5; i++ {
		arr[i] = &TrainingRecord{
			TrainingID: fmt.Sprintf("Training-%d-%d", i, time.Now().Unix()),
			UserID:     "test-user1",
			ModelDefinition: &grpc_trainer_v2.ModelDefinition{
				Name: "foo-name",
			},
		}
	}
	for i := 5; i < 10; i++ {
		arr[i] = &TrainingRecord{
			TrainingID: fmt.Sprintf("Training-%d-%d", i, time.Now().Unix()),
			UserID:     "test-user2",
		}
	}
	return arr
}

func TestMongoRepository(t *testing.T) {
	logEntry().Debugf("Using Mongo address: %s", mongoAddress)

	collectionName := fmt.Sprintf("itest_trainer_repo_%d", time.Now().Unix())
	r, err := newTrainingsRepository(mongoAddress, mongoDatabase, mongoUsername, mongoPassword,
		mongoCertLocation, collectionName)
	assert.NoError(t, err)
	defer r.Close()

	// create 10 Training instances and Store()
	trainings := createTrainings()
	for i := 0; i < len(trainings); i++ {
		assert.NoError(t, r.Store(trainings[i]))
	}

	// Find()
	found, err := r.Find("non_existing_id") // expect both to be nil. a non-existing item is not an error
	assert.Nil(t, found)
	if assert.Error(t, err) {
		assert.Equal(t, err, mgo.ErrNotFound)
	}

	found, err = r.Find(trainings[0].TrainingID)
	assert.NotNil(t, found)
	assert.NoError(t, err)
	assert.EqualValues(t, trainings[0].TrainingID, found.TrainingID)
	assert.EqualValues(t, trainings[0].UserID, found.UserID)
	assert.EqualValues(t, trainings[0].ModelDefinition.Name, found.ModelDefinition.Name)

	records, err := r.FindAll("test-user1", "", "")
	assert.NotNil(t, records)
	assert.NoError(t, err)
	assert.EqualValues(t, 5, 5)

	records, err = r.FindAll("test-user2", "", "")
	assert.NotNil(t, records)
	assert.NoError(t, err)
	assert.EqualValues(t, 5, 5)

	// Testing updating a training record
	tr, err := r.Find(trainings[0].TrainingID)
	assert.NotNil(t, found)
	assert.NoError(t, err)

	tr.ModelDefinition.Name = "bar-name"
	err = r.Store(tr)
	assert.NoError(t, err)

	tr2, err := r.Find(trainings[0].TrainingID)
	assert.NotNil(t, tr2)
	assert.NoError(t, err)
	assert.EqualValues(t, "bar-name", tr2.ModelDefinition.Name)

	// manually connect to Mongo to check state of soft-deleted records
	sess, _ := ConnectMongo(mongoAddress, mongoDatabase, mongoUsername, mongoPassword, mongoCertLocation)
	coll := sess.DB(mongoDatabase).C(collectionName)
	defer sess.Close()

	// Delete()
	for i := 0; i < len(trainings); i++ {
		trainingID := trainings[i].TrainingID
		assert.NoError(t, r.Delete(trainingID))

		// validate deletion
		found, err = r.Find(trainingID)
		assert.Nil(t, found)
		if assert.Error(t, err) {
			assert.Equal(t, err, mgo.ErrNotFound)
		}

		// validate that the records still exist in Mongo (soft-delete)
		count, _ := coll.Find(bson.M{"training_id": trainingID}).Count()
		assert.True(t, count == 1)
	}

}

//assert that FindCurrentlyRunningTrainings gets trainings sorted by the order of creation in descending manner
func TestMongoRespositoryFindCurrentlyRunningTrainings(t *testing.T) {

	var createTrainingsWithGPUs = func(count int) []*TrainingRecord {
		arr := make([]*TrainingRecord, count)
		for i := 0; i < count; i++ {
			arr[i] = &TrainingRecord{
				TrainingID: fmt.Sprintf("%d", time.Now().UnixNano()),
				UserID:     "TestMongoRespositoryFindCurrentlyRunningTrainings",
				ModelDefinition: &grpc_trainer_v2.ModelDefinition{
					Name: "foo-name",
				},
				Training: &grpc_trainer_v2.Training{
					Resources: &grpc_trainer_v2.ResourceRequirements{
						Gpus: 1,
						Cpus: 1,
					},
				},
				TrainingStatus: &grpc_trainer_v2.TrainingStatus{
					Status: grpc_trainer_v2.Status_PENDING,
				},
			}
		}
		return arr
	}

	logEntry().Debugf("Using Mongo address: %s", mongoAddress)

	r, err := newTrainingsRepository(mongoAddress, mongoDatabase, mongoUsername, mongoPassword,
		mongoCertLocation, fmt.Sprintf("itest_trainer_repo_%d", time.Now().Unix()))
	assert.NoError(t, err)
	defer r.Close()

	// create 10 Training instances and Store()
	trainings := createTrainingsWithGPUs(10)
	for i := 0; i < len(trainings); i++ {
		assert.NoError(t, r.Store(trainings[i]))
	}

	currentTrainings, err := r.FindCurrentlyRunningTrainings(len(trainings))
	assert.NoError(t, err)

	//assert that the trainings are sorted in the descending order of time of creation
	lastTrainingID := time.Now().UnixNano()

	for _, record := range currentTrainings {
		assert.True(t, cast.ToInt64(record.TrainingID) < lastTrainingID, "Last training id was %v and comparing it with current id %v", record.TrainingID, lastTrainingID)
		lastTrainingID = cast.ToInt64(record.TrainingID)
		//assert that you can get the gpu and cpu counts and those fields are present in the result
		assert.True(t, record.Training.Resources.Gpus > 0)
		assert.True(t, record.Training.Resources.Cpus > 0)
		//assert that training status field is present on all the records
		assert.NotNil(t, record.TrainingStatus.Status)
	}

	// Delete()
	for i := 0; i < len(trainings); i++ {
		assert.NoError(t, r.Delete(trainings[i].TrainingID))
	}

}

// FIXME MLSS Change: test TestFindAllByUserIdAndNamespaceList
func TestFindAllByUserIdAndNamespaceList(t *testing.T) {
	//logr := logger.LocLogger(log.StandardLogger().WithField("module", "trainingRepository"))trainingRepository
	//var tr []*TrainingRecord
	//sess := r.session.Clone()
	//defer sess.Close()
	//create query object
	s := "zhou"
	user := &s
	strings := []string{"1", "2"}
	namespaces := &strings
	//m := &bson.M{}
	//if *user != "" {
	//	(*m)["user_id"] = *user
	//}
	//if len(*namespaces) > 0 {
	//	var n []bson.M
	//	for _, v := range *namespaces {
	//		n = append(n, bson.M{"namespace": v})
	//	}
	//	(*m)["$or"] = n
	//}
	//logr.Infof("FindAllByUserIdAndNamespaceList query info: %v", *m)

	//m := &bson.M{}
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
	fmt.Println(len(n))
	//logr.Infof("FindAllByUserIdAndNamespaceList query info: %v", *m)
	if len(n) == 0 {
		//return nil, errors.New("there is nothing to query")
	}
}
