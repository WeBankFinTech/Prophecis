package mongo

import (
	mgo "gopkg.in/mgo.v2"
	"os"
	"time"
	"webank/AIDE/notebook-server/pkg/commons/logger"
)

//var Session *mgo.Session
var MongoDatabase string

const (
	MONGO_ADDRESS                = "MONGO_ADDRESS"
	MONGO_USER_NAME              = "MONGO_USER_NAME"
	MONGO_PASSWORD               = "MONGO_PASSWORD"
	MONGO_DATABASE               = "MONGO_DATABASE"
	NoteBookCollection           = "notebooks_collection"
	MONGO_Authentication_Database = "MONGO_Authentication_Database"
)

func GetMongoSession() (*mgo.Session, error) {
	address := os.Getenv(MONGO_ADDRESS)
	username := os.Getenv(MONGO_USER_NAME)
	passwd := os.Getenv(MONGO_PASSWORD)
	MongoDatabase := os.Getenv(MONGO_DATABASE)
	if MongoDatabase == "" {
		MongoDatabase = "mlss"
	}
	authenticationDatabase := os.Getenv(MONGO_Authentication_Database)
	if authenticationDatabase == "" {
		authenticationDatabase = "admin"
	}
	logger.Logger().Debugf("mongo address: %q, user name: %q, password: %q, database: %q, authenticationDatabase: %q",
		address, username, passwd, MongoDatabase, authenticationDatabase)
	dialInfo := mgo.DialInfo{
		Addrs:          []string{address},
		Direct:         false,
		Timeout:        time.Second * 1,
		FailFast:       false,
		Database:       MongoDatabase,
		ReplicaSetName: "",
		//Source:         "admin",
		Source:      authenticationDatabase,
		Service:     "",
		ServiceHost: "",
		Mechanism:   "",
		Username:    username,
		Password:    passwd,
		PoolLimit:   4096,
		DialServer:  nil,
		Dial:        nil,
	}

	return mgo.DialWithInfo(&dialInfo)
}

func Init() {
	session, err := GetMongoSession()
	defer session.Close()
	if err != nil {
		logger.Logger().Errorf("mongo dial fail, err: %v\n", err)
		return
	}
	logger.Logger().Infoln("mongo session init success......")

	initMongo(MongoDatabase, session)
	syncNotebookToMongo(session)
}

func initMongo(db string, session *mgo.Session) {
	//创建集合
	sessionClone := session.Clone()
	defer sessionClone.Close()
	collectionNames, err := sessionClone.DB(db).CollectionNames()
	if err != nil {
		logger.Logger().Errorf("mongo get collection fail, error: %v\n", err)
		return
	}
	isExist := false
	for _, v := range collectionNames {
		logger.Logger().Debugf("mongo collection name: %q", v)
		if v == NoteBookCollection {
			isExist = true
			break
		}
	}

	if !isExist {
		//err = session.DB(db).C("notebooks_collection").EnsureIndexKey("namespace", "name")
		index := mgo.Index{
			Key:              []string{"namespace", "name", "enable_flag"},
			Unique:           false,
			DropDups:         false,
			Background:       true,
			Sparse:           true,
			ExpireAfter:      0,
			Name:             "namespace-name-enable_flag",
			Min:              0,
			Max:              0,
			Minf:             0,
			Maxf:             0,
			BucketSize:       0,
			Bits:             0,
			DefaultLanguage:  "",
			LanguageOverride: "",
			Weights:          nil,
			Collation:        nil,
		}
		err = sessionClone.DB(db).C(NoteBookCollection).EnsureIndex(index)
		if err != nil {
			logger.Logger().Errorf("mongo create index fail: %v", err)
			return
		}
		idIndex := mgo.Index{
			Key:              []string{"id"},
			Unique:           true,
			DropDups:         false,
			Background:       true,
			Sparse:           true,
			ExpireAfter:      0,
			Name:             "id",
			Min:              0,
			Max:              0,
			Minf:             0,
			Maxf:             0,
			BucketSize:       0,
			Bits:             0,
			DefaultLanguage:  "",
			LanguageOverride: "",
			Weights:          nil,
			Collation:        nil,
		}
		err = sessionClone.DB(db).C(NoteBookCollection).EnsureIndex(idIndex)
		if err != nil {
			logger.Logger().Errorf("mongo create id index fail: %v", err)
			return
		}
	}

	logger.Logger().Infoln("mongo init success......")
}

//const (
//	sslSuffix = "?ssl=true"
//)
//
//func ConnectMongo(mongoURI string, database string, username string, password string, cert string) (*mgo.Session, error) {
//	uri := strings.TrimSuffix(mongoURI, sslSuffix)
//	dialInfo, err := mgo.ParseURL(uri)
//	if err != nil {
//		log.WithError(err).Errorf("Cannot parse Mongo Connection URI")
//		return nil, err
//	}
//	dialInfo.FailFast = true
//	dialInfo.Timeout = 10 * time.Second
//
//	// only do ssl if we have the suffix
//	if strings.HasSuffix(mongoURI, sslSuffix) {
//		log.Debugf("Using TLS for mongo ...")
//		tlsConfig := &tls.Config{}
//		roots := x509.NewCertPool()
//		if ca, err := ioutil.ReadFile(cert); err == nil {
//			roots.AppendCertsFromPEM(ca)
//		}
//		tlsConfig.RootCAs = roots
//		tlsConfig.InsecureSkipVerify = false
//		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
//			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
//			return conn, err
//		}
//	}
//
//	// in case the username/password are not part of the URL string
//	if username != "" && password != "" {
//		dialInfo.Username = username
//		dialInfo.Password = password
//	}
//
//	session, err := mgo.DialWithInfo(dialInfo)
//
//	if database == "" {
//		database = dialInfo.Database
//	}
//
//	if err != nil {
//		msg := fmt.Sprintf("Cannot connect to MongoDB at %s, db %s", mongoURI, database)
//		log.WithError(err).Errorf(msg)
//		return nil, err
//	}
//
//	return session, nil
//}
