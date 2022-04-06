package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"webank/AIDE/notebook-server/pkg/commons/logger"
	"webank/AIDE/notebook-server/pkg/commons/utils"
	"webank/AIDE/notebook-server/pkg/models"
	"encoding/json"
)

func syncNotebookToMongo(session *mgo.Session) {
	ncClient := utils.GetNBCClient()
	ns := ncClient.GetNamespaces()
	if len(ns) > 0 {
		for _, namespace := range ns {
			listFromK8s, err := ncClient.GetNotebooks(namespace, "", "")
			if err != nil {
				logger.Logger().Errorf("fail to get notebook in namespace %q, err: %v", namespace, err)
				continue
			}
			configFromK8s, err := ncClient.GetNBConfigMaps(namespace)
			if err != nil {
				logger.Logger().Errorf("fail to get notebook configmap in namespace %q, err: %v", namespace, err)
				continue
			}

			notebooks, err := parseNotebooksFromK8s(listFromK8s)
			if err != nil {
				logger.Logger().Errorf("fail to parse notebook, listFromK8s: %v, err: %v", listFromK8s, err)
				continue
			}
			notebookInMongos, err := utils.ParseToNotebookInMongo(*notebooks, *configFromK8s)
			if err != nil {
				logger.Logger().Errorf("fail to parse notebook to mongo, notebooks: %v, "+
					"configFromK8s: %v, err: %v", *notebooks, *configFromK8s, err)
				continue
			}
			if len(notebookInMongos) > 0 {
				// sync to mongo
				sessionClone := session.Clone()
				defer sessionClone.Close()
				for _, nb := range notebookInMongos {
					count, err := sessionClone.DB(MongoDatabase).C(NoteBookCollection).Find(bson.M{"namespace": nb.Namespace,
						"name": nb.Name, "enable_flag": nb.EnableFlag}).Count()
					if err != nil {
						logger.Logger().Errorf("mongo fail to get notebook, namespace: %q, "+
							"name: %q, enable_flag: %v", nb.Namespace, nb.Name, nb.EnableFlag)
						continue
					}
					if count > 0 {
						logger.Logger().Errorf("mongo notebook has exist in mongo, namespace: %q, "+
							"name: %q, enable_flag: %v, notebook: %+v", nb.Namespace, nb.Name, nb.EnableFlag, nb)
						continue
					} else {
						err := sessionClone.DB(MongoDatabase).C(NoteBookCollection).Insert(nb)
						if err != nil {
							logger.Logger().Errorf("mongo fail to insert notebook, notebook: %+v", nb)
							continue
						}
					}
				}
			}
		}
	}
	logger.Logger().Info("mongo sync notebook success......")
}

func parseNotebooksFromK8s(listFromK8s interface{}) (*models.NotebookFromK8s, error) {
	var notebookFromK8s models.NotebookFromK8s
	bytes, err := json.Marshal(listFromK8s)
	if err != nil {
		logger.Logger().Errorf("Marshal failed. %v", err.Error())
		return nil, err
	}
	err = json.Unmarshal(bytes, &notebookFromK8s)
	if err != nil {
		logger.Logger().Errorf("parse []models.K8sNotebook Unmarshal failed. %v", err.Error())
		return nil, err
	}
	return &notebookFromK8s, nil
}
