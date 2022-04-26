package datasource

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/schema"
	"webank/DI/commons/logger"
)
import "gorm.io/gorm"
import _ "gorm.io/driver/mysql"

var DB *gorm.DB

//var MlssDS DataSource
//var TokenDS DataSource

const (
	tokendb_username = "spring.tokendb.username"
	tokendb_password = "spring.tokendb.password"
	tokendb_ip       = "spring.tokendb.ip"
	tokendb_port     = "spring.tokendb.port"
	tokendb_db       = "spring.tokendb.db"
	token_db_name    = "WB_USER_TOKEN"
	datasource_user_name = "datasource.userName"
	datasource_user_pwd = "datasource.userPwd"
	datasource_encrypt_pwd = "datasource.encryptPwd"
	datasource_url = "datasource.url"
	datasource_port= "datasource.port"
	datasource_db= "datasource.db"
	datasource_priv_key= "datasource.privKey"

	WbUserToken = "WB_USER_TOKEN"
)


func InitDS(encodeFlag bool) {
	log := logger.GetLogger()
	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return "t_" + defaultTableName
	//}
	var password string
	if encodeFlag {
		decryptPwd, err := RsaDecrypt()
		if err != nil {
			logger.GetLogger().Error("Datasource rsa decrypt error, ", err.Error())
			return
		}
		logger.GetLogger().Info("Rsa decrypt password：",string(decryptPwd))
		password = string(decryptPwd)
	}else{
		password = viper.GetString(datasource_user_pwd)
	}
	dbUrl := fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		viper.GetString(datasource_user_name),
		password,
		viper.GetString(datasource_url),
		viper.GetString(datasource_port),
		viper.GetString(datasource_db),
	)

	log.Infof("mlss-db url: %v", dbUrl)

	//db, err := gorm.Open("mysql", dbUrl)
	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "t_", SingularTable: true}})
	if err != nil {
		log.Fatalf("failed to connect database, reason: %v", err.Error())
	}
	//db.SingularTable(true)

	//debug
	db = db.Debug()

	DB = db

	//TokenDS = initTokenDB()
}

func GetDB() *gorm.DB {
	return DB
}

func initTokenDB() DataSource {
	log := logger.GetLogger()

	dbUrl := fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		//viper.GetString(tokendb_username),
		//viper.GetString(tokendb_password),
		//viper.GetString(tokendb_ip),
		//viper.GetString(tokendb_port),
		//viper.GetString(tokendb_db),
		"",
		"",
		"",
		"3306",
		"cnc_hdfs_privs",
	)
	log.Infof("tokendb url: %v", dbUrl)

	//db, err := gorm.Open("mysql", dbUrl)
	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "t_", SingularTable: true}})
	if err != nil {
		log.Fatalf("failed to connect tokenDB, reason: %v", err.Error())
	}
	//db.SingularTable(true)

	source := DataSource{
		DB: db,
	}

	return source

}

type DataSource struct {
	DB *gorm.DB
}

func (ds *DataSource) GetDB() *gorm.DB {
	return ds.DB
}

func RsaDecrypt() ([]byte, error) {
	encryptPwd := viper.GetString(datasource_encrypt_pwd)
	if len(encryptPwd) < 8 {
		return []byte(viper.GetString(datasource_encrypt_pwd)), errors.New("encryptPwd length error")
	}
	decode, err := hex.DecodeString(encryptPwd[8:])
	if err != nil {
		logger.GetLogger().Error("Datasource rsa decode error, ", err.Error())
		return []byte(""), nil
	}
	privKey := "-----BEGIN RSA PRIVATE KEY-----\n" + viper.GetString(datasource_priv_key) + "\n-----END RSA PRIVATE KEY-----"
	block, _ := pem.Decode([]byte(privKey))
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv.(*rsa.PrivateKey), decode)
}

//func (ds *DataSource) GetDBByDBS(txs []*gorm.DB) *gorm.DB {
//	var db *gorm.DB
//	if len(txs) > 0 {
//		db = txs[0]
//	} else {
//		db = ds.GetDB().New()
//		defer db.Close()
//
//	}
//	return db
//}
