### Prophecis编译文档

[TOC]

#### 1、Go Build

```shell
# prophecis-cc-apiserver
go build -v -installsuffix cgo -ldflags '-w' -o mlss-controlcenter-go cmd/mlss-controlcenter-go/main.go

# prophecis-cc-apigateway
go build -v -a -installsuffix cgo -ldflags '-w' -o mlss-apigateway cmd/mlss-apigateway/main.go

# prophecis-mllabis
go build -v -o build/notebook-server cmd/jupyter-server/main.go
```

#### 2、Docker Build

```shell
cd $WORKSPACE/notebook-server
docker build -t wedatasphere/prophecis:mllabis-v0.1.0 .


cd $WORKSPACE/build/mlss-apiserver/
docker build  wedatasphere/prophecis:ccapiserver-v0.1.0 .

cd $WORKSPACE/build/mlss-apigateway/.
docker build wedatasphere/prophecis:ccapigateway-v0.1.0
```

#### 3、其他

- swagger命令

  Prophecis的Rest API是通过go-swagger进行生成的，可以通过以下命令生成相应的API

```shell
# CC
swagger generate server -m ../../models -f ../swagger/swagger.yaml -t ./ -A  mlss_cc -s ./ —exclude-main
```



- 关于gateway，prophecis gateway使用caddy对客户端的请求进行拦截和转发，为此定制对应的自定义插件，因此在编译时需要改动下caddy plugin的代码：

  - 在https://github.com/caddyserver/caddy/blob/d3860f95f59b5f18e14ddf3d67b4c44dbbfdb847/caddyhttp/httpserver/plugin.go#L314-L355中的该段代码中，增加auth_request插件：

    ```
    var directives = []string{
    	// primitive actions that set up the fundamental vitals of each config
    	"root",
    	"tls",
    	"bind",
    
    	// services/utilities, or other directives that don't necessarily inject handlers
    	"startup",
    	"shutdown",
    	"realip", // github.com/captncraig/caddy-realip
    	"git",    // github.com/abiosoft/caddy-git
    
    	// directives that add middleware to the stack
    	"log",
    	"gzip",
    	"errors",
    	"ipfilter", // github.com/pyed/ipfilter
    	"search",   // github.com/pedronasser/caddy-search
    	"header",
    	"cors", // github.com/captncraig/cors/caddy
    	"rewrite",
    	"redir",
    	"ext",
    	"mime",
    	"basicauth",
    	"jwt",    // github.com/BTBurke/caddy-jwt
    	"jsonp",  // github.com/pschlump/caddy-jsonp
    	"upload", // blitznote.com/src/caddy.upload
    	"internal",
    	"proxy",
    	"fastcgi",
    	"websocket",
    	"markdown",
    	"templates",
    	"browse",
    	"hugo",       // github.com/hacdias/caddy-hugo
    	"mailout",    // github.com/SchumacherFM/mailout
    	"prometheus", // github.com/miekg/caddy-prometheus
    }
    ```