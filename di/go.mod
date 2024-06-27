module webank/DI

go 1.13

require (
	github.com/IBM-Bluemix/bluemix-cli-sdk v0.6.7
	github.com/IBM-Cloud/ibm-cloud-cli-sdk v0.6.7 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d
	github.com/aws/aws-sdk-go v1.42.50
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/dre1080/recover v0.0.0-20150930082637-1c296bbb3227
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-kit/kit v0.10.0
	github.com/go-openapi/errors v0.20.2
	github.com/go-openapi/loads v0.21.0
	github.com/go-openapi/runtime v0.21.1
	github.com/go-openapi/spec v0.20.4
	github.com/go-openapi/strfmt v0.21.1
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.3
	github.com/golang/glog v1.0.0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/kubeflow/pipelines v0.0.0-20231208180049-c5658f09ec38
	github.com/mholt/archiver/v3 v3.3.0
	github.com/minio/minio-go/v6 v6.0.57
	github.com/modern-go/reflect2 v1.0.2
	github.com/ncw/swift v1.0.50
	github.com/nicksnyder/go-i18n v1.10.1 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/common v0.32.1
	github.com/sirupsen/logrus v1.8.1
	github.com/sony/gobreaker v0.4.1
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.8.1
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf // indirect
	github.com/thoas/go-funk v0.6.0
	github.com/tylerb/graceful v1.2.15
	github.com/urfave/cli v1.22.4
	github.com/ventu-io/go-shortid v0.0.0-20171029131806-771a37caa5cf
	golang.org/x/net v0.21.0
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.30.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/olivere/elastic.v5 v5.0.85
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.0.2
	gorm.io/gorm v1.20.2
	k8s.io/api v0.24.3
	k8s.io/apimachinery v0.24.3
	k8s.io/client-go v0.24.3
	k8s.io/kube-openapi v0.0.0-20220627174259-011e075b9cb8
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)

replace (
	//由于引入kfp后go-openapi/xxx的版本升级了，但是我们的restapi是用的旧版的go-openapi生成的，有很多函数接口发生变化了
	github.com/go-openapi/analysis => github.com/go-openapi/analysis v0.19.7
	github.com/go-openapi/errors => github.com/go-openapi/errors v0.19.4
	github.com/go-openapi/loads => github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/runtime => github.com/go-openapi/runtime v0.19.15
	github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.7
	github.com/go-openapi/strfmt => github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag => github.com/go-openapi/swag v0.19.8
	github.com/go-openapi/validate => github.com/go-openapi/validate v0.19.7
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1
	github.com/googleapis/gnostic/OpenAPIv2 => github.com/googleapis/gnostic/OpenAPIv2 v0.3.1
	github.com/kubeflow/pipelines v0.0.0-20231208180049-c5658f09ec38 => localhost/MLSS/MLSS-TRAINERDIPIPELINE v0.0.0-20240527134309-ded58d882198
	github.com/nicksnyder/go-i18n/i18n => github.com/nicksnyder/go-i18n/i18n v1.10.1
	google.golang.org/grpc => google.golang.org/grpc v1.29.0
	k8s.io/api => github.com/kubernetes/api v0.16.8
	k8s.io/client-go => github.com/kubernetes/client-go v0.16.8
	k8s.io/kubernetes => k8s.io/kubernetes v1.11.1
)
