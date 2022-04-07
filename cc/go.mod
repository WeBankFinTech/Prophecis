module mlss-controlcenter-go

go 1.13

require (
	github.com/IBM/FfDL v0.1.1-rc
	github.com/caddyserver/caddy v1.0.4
	github.com/dlclark/regexp2 v1.2.0
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/gin-contrib/location v0.0.1
	github.com/gin-gonic/gin v1.5.0
	github.com/go-ldap/ldap/v3 v3.2.4
	github.com/go-openapi/errors v0.19.3
	github.com/go-openapi/loads v0.19.4
	github.com/go-openapi/runtime v0.19.11
	github.com/go-openapi/spec v0.19.6
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.19.7
	github.com/go-openapi/validate v0.19.5
	github.com/google/uuid v1.1.1
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/jinzhu/copier v0.3.5
	github.com/jinzhu/gorm v1.9.12
	github.com/kubernetes-client/go v0.0.0-20200222171647-9dac5e4c5400
	github.com/lestrrat-go/file-rotatelogs v2.3.0+incompatible
	github.com/lestrrat-go/strftime v1.0.1 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/common v0.10.0
	github.com/seldonio/seldon-core/operator v0.0.0-20210122151420-b881f99ebc1f
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.0
	github.com/tebeka/strftime v0.1.3 // indirect
	golang.org/x/net v0.0.0-20201202161906-c7110b5ffcbb
	gopkg.in/cas.v2 v2.2.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73 // indirect
)

replace k8s.io/api v0.18.8 => k8s.io/api v0.18.6

replace k8s.io/apimachinery v0.18.8 => k8s.io/apimachinery v0.18.6

replace k8s.io/client-go v12.0.0+incompatible => k8s.io/client-go v0.18.6
