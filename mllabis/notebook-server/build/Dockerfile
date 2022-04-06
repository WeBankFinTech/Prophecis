#FROM ubuntu:latest
FROM golang-dlv:latest
WORKDIR /
RUN echo "Asia/shanghai" > /etc/timezone
COPY notebook-server .
ENTRYPOINT ["/notebook-server"]