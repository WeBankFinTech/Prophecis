module webank/DI

go 1.13

require (
	github.com/IBM-Bluemix/bluemix-cli-sdk v0.6.7
	github.com/IBM-Cloud/ibm-cloud-cli-sdk v0.6.7 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d
	github.com/aws/aws-sdk-go v1.29.11
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/coreos/etcd v3.3.10+incompatible
	github.com/dre1080/recover v0.0.0-20150930082637-1c296bbb3227
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-kit/kit v0.10.0
	github.com/go-openapi/errors v0.19.4
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/runtime v0.19.15
	github.com/go-openapi/spec v0.19.7
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.8
	github.com/go-openapi/validate v0.19.7
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.4.1
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gorilla/websocket v1.4.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/mholt/archiver/v3 v3.3.0
	github.com/minio/minio-go/v6 v6.0.57
	github.com/modern-go/reflect2 v1.0.1
	github.com/ncw/swift v1.0.50
	github.com/nicksnyder/go-i18n v1.10.1 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1
	github.com/prometheus/common v0.9.1
	github.com/sirupsen/logrus v1.5.0
	github.com/sony/gobreaker v0.4.1
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v0.0.7
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf // indirect
	github.com/thoas/go-funk v0.6.0
	github.com/tylerb/graceful v1.2.15
	github.com/urfave/cli v1.22.4
	github.com/ventu-io/go-shortid v0.0.0-20171029131806-771a37caa5cf
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	google.golang.org/grpc v1.28.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/olivere/elastic.v5 v5.0.85
	gopkg.in/yaml.v2 v2.2.8
	gorm.io/driver/mysql v1.0.2
	gorm.io/gorm v1.20.2
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.0.0-00010101000000-000000000000
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
)

replace k8s.io/client-go => github.com/kubernetes/client-go v0.16.8

replace github.com/googleapis/gnostic/OpenAPIv2 => github.com/googleapis/gnostic/OpenAPIv2 v0.3.1

replace github.com/nicksnyder/go-i18n/i18n => github.com/nicksnyder/go-i18n/i18n v1.10.1
