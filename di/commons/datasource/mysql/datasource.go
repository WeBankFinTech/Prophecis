package datasource

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"webank/DI/commons/logger"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	_ "gorm.io/driver/mysql"
)

var DB *gorm.DB

//var MlssDS DataSource
//var TokenDS DataSource

const (
	token_db_name          = "WB_USER_TOKEN"
	datasource_user_name   = "datasource.userName"
	datasource_user_pwd    = "datasource.userPwd"
	datasource_encrypt_pwd = "password"
	datasource_url         = "datasource.url"
	datasource_port        = "datasource.port"
	datasource_db          = "datasource.db"
	datasource_priv_key    = "datasource.privKey"

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
		logger.GetLogger().Info("Rsa decrypt password：", string(decryptPwd))
		password = string(decryptPwd)
	} else {
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

	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "t_", SingularTable: true}})
	if err != nil {
		log.Fatalf("failed to connect database, reason: %v", err.Error())
	}

	// db = db.Debug()

	DB = db

}

func GetDB() *gorm.DB {
	return DB
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
	privKey := "***REMOVED***----\n" + viper.GetString(datasource_priv_key) + "\n-----END RSA PRIVATE KEY-----"
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
