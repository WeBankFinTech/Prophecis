#FROM ubuntu:latest
#FROM golang:1.13.15
FROM golang-dlv:latest
WORKDIR /
COPY mf-server .
#COPY localhost/MLSS/MLSS-MF/-client-app-1.4.0-fat.jar .
RUN echo "Asia/shanghai" > /etc/timezone
#ENV GOPROXY http://localhost
#ENV GO111MODULE=on
#GOPROXY="http://127.0.0.1"
#RUN go get -u github.com/go-delve/delve/cmd/dlv
EXPOSE 40000
#ENTRYPOINT ["/mf-server"]
ENTRYPOINT [ "/go/bin/dlv" , "--listen=:40000" , "--headless=true" , "--api-version=2" , "exec" , "/mf-server" ]