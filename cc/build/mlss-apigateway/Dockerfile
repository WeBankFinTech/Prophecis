#FROM dev.sf.dockerhub.stgwebank/webank/mlss-go-cc:golang-1.12
#FROM uat.sf.dockerhub.stgwebank/webank/mlss-di:ubuntu-latest
#FROM ubuntu:latest
FROM golang-dlv:latest

COPY mlss-apigateway /apigateway/
WORKDIR /apigateway

CMD ["./mlss-apigateway"]

EXPOSE 8080