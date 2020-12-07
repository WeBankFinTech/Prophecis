all: build push
push:
	docker push buoylee/cc
build: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/MLSS-ControlCenter-go/main.go
	docker build -t buoylee/cc .
clean:
	rm -f main

gen-swagger:
	swagger generate server -m ../models -f pkg/restapi/swagger/swagger.yaml -t pkg/restapi --exclude-main