#FROM dev.sf.dockerhub.stgwebank/webank/mlss-go-cc:golang-1.12
#FROM uat.sf.dockerhub.stgwebank/webank/mlss-di:ubuntu-latest
FROM golang-dlv:latest
COPY mlss-controlcenter-go /cc/
WORKDIR /cc

CMD ["./mlss-controlcenter-go"]

EXPOSE 6666