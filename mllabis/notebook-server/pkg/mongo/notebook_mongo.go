package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"webank/AIDE/notebook-server/pkg/commons/logger"
	"webank/AIDE/notebook-server/pkg/commons/utils"
)

func SearchNoteBook(query interface{}) ([]utils.NotebookInMongo, error) {
	session, err := GetMongoSession()
	defer session.Close()
	if err != nil {
		return nil, err
	}
	result := []utils.NotebookInMongo{}
	err = session.DB(MongoDatabase).C(NoteBookCollection).Find(query).All(&result)
	return result, err
}

func CountNoteBook(query interface{}) (int, error) {
	session, err := GetMongoSession()
	defer session.Close()
	if err != nil {
		return 0, err
	}

	count, err := session.DB(MongoDatabase).C(NoteBookCollection).Find(query).Count()
	return count, err
}

func InsertNoteBook(nb utils.NotebookInMongo) error {
	session, err := GetMongoSession()
	defer session.Close()
	if err != nil {
		return err
	}

	err = session.DB(MongoDatabase).C(NoteBookCollection).Insert(nb)
	return err
}

func ListNootBook(namespace, role, userId string) ([]utils.NotebookInMongo, error) {
	notebooks := make([]utils.NotebookInMongo, 0)
	session, err := GetMongoSession()
	defer session.Close()
	if err != nil {
		return nil, err
	}

	query := bson.M{"enable_flag": 1}
	if namespace != "" && namespace != "null"{
		query["namespace"] = namespace
	}
	if role == "GU" && userId != "" {
		query["user"] = userId
	}

	//err := session.DB(MongoDatabase).C(NoteBookCollection).Find(query).Sort("id").Skip(offset).Limit(limit).All(&result)
	err = session.DB(MongoDatabase).C(NoteBookCollection).Find(query).All(&notebooks)
	return notebooks, err
}

func ListNamespaceUserNotebook(namespace, userId string) ([]utils.NotebookInMongo, error) {
	notebooks := make([]utils.NotebookInMongo, 0)
	session, err := GetMongoSession()
	defer session.Close()
	if err != nil {
		return nil, err
	}

	query := bson.M{"enable_flag": 1}
	if namespace != "" && namespace != "null" {
		query["namespace"] = namespace
	}
	if userId != "" && userId != "null" {
		query["user"] = userId
	}

	logger.Logger().Debugf("ListNamespaceUserNotebook query: %+v", query)

	//err := session.DB(MongoDatabase).C(NoteBookCollection).Find(query).Sort("id").Skip(offset).Limit(limit).All(&result)
	err = session.DB(MongoDatabase).C(NoteBookCollection).Find(query).All(&notebooks)
	return notebooks, err
}


func ListUserNotebook( userId string) ([]utils.NotebookInMongo, error) {
	notebooks := make([]utils.NotebookInMongo, 0)
	session, err := GetMongoSession()
	defer session.Close()
	if err != nil {
		return nil, err
	}

	query := bson.M{"enable_flag": 1}
	if userId != "" && userId != "null" {
		query["user"] = userId
	} else {

	}

	logger.Logger().Debugf("ListNamespaceUserNotebook query: %+v", query)

	//err := session.DB(MongoDatabase).C(NoteBookCollection).Find(query).Sort("id").Skip(offset).Limit(limit).All(&result)
	err = session.DB(MongoDatabase).C(NoteBookCollection).Find(query).All(&notebooks)
	return notebooks, err
}

func GetNotebookById(id string) (utils.NotebookInMongo, error) {
	notebook := utils.NotebookInMongo{}
	session, err := GetMongoSession()
	defer session.Close()
	if err != nil {
		return notebook, err
	}

	query := bson.M{"id": id, "enable_flag": 1}

	//err := session.DB(MongoDatabase).C(NoteBookCollection).Find(query).Sort("id").Skip(offset).Limit(limit).All(&result)
	err = session.DB(MongoDatabase).C(NoteBookCollection).Find(query).One(&notebook)
	return notebook, err
}

func DeleteNotebookById(id string) error {
	session, err := GetMongoSession()
	defer session.Close()
	if err != nil {
		return  err
	}

	query := bson.M{"id": id, "enable_flag": 1}
	upadteMap := bson.M{"$set": bson.M{"enable_flag": 0}}

	//err := session.DB(MongoDatabase).C(NoteBookCollection).Find(query).Sort("id").Skip(offset).Limit(limit).All(&result)
	err = session.DB(MongoDatabase).C(NoteBookCollection).Update(query, upadteMap)
	return err
}

func UpdateNotebookById(id string, m bson.M) error {
	session, err := GetMongoSession()
	defer session.Close()
	if err != nil {
		return  err
	}

	query := bson.M{"id": id, "enable_flag": 1}
	upadteMap := bson.M{"$set": m}

	//err := session.DB(MongoDatabase).C(NoteBookCollection).Find(query).Sort("id").Skip(offset).Limit(limit).All(&result)
	err = session.DB(MongoDatabase).C(NoteBookCollection).Update(query, upadteMap)
	return err
}

//func getNoteBookByClusterName(clusterName string, res []*models.Notebook) []*models.Notebook {
//	clusterResult := make([]*models.Notebook, 0)
//	if len(res) > 0 {
//		if clusterName != "" {
//			for _, v := range res {
//				if clusterName == utils.BDP || clusterName == utils.BDAP {
//					if stringsUtil.Split(v.Namespace, "-")[2] == clusterName && !stringsUtil.Contains(stringsUtil.Split(v.Namespace, "-")[6], "safe") {
//						clusterResult = append(clusterResult, v)
//					}
//				}
//				if clusterName == utils.BDAPSAFE {
//					if stringsUtil.Split(v.Namespace, "-")[2] == utils.BDAP && stringsUtil.Contains(stringsUtil.Split(v.Namespace, "-")[6], "safe") {
//						clusterResult = append(clusterResult, v)
//					}
//				}
//			}
//			//res = clusterResult
//		} else {
//			clusterResult = res
//		}
//	}
//	return clusterResult
//}
