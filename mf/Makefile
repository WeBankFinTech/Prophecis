registry = wedatasphere
repo = prophecis
version = v0.3.0
tag = mf-server-$(version)

build-bin:
	CGO_ENABLED=0 GOSUMDB=off GOOS=linux GOARCH=amd64 /data/Install/goInstall/go/bin/go build -v -o build/mf-server cmd/mf-server/main.go
build-image: build-bin
	docker build -t $(registry)/$(repo):$(tag) -f build/Dockerfile build
push-image:
	docker push $(registry)/$(repo):$(tag)
build: build-image push-image
test:
	echo $(registry)/$(repo):$(tag)