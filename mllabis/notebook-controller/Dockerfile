# Build the manager binary
# Copy the controller-manager into a thin image
FROM ubuntu:latest
WORKDIR /
COPY  ./manager .
ENTRYPOINT ["/manager"]