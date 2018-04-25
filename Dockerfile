FROM golang:alpine as builder
COPY . /go/src/code.cloudfoundry.org/slack-attachment-resource
RUN go build -o /opt/resource/check code.cloudfoundry.org/slack-attachment-resource/check
RUN go build -o /opt/resource/in code.cloudfoundry.org/slack-attachment-resource/in
RUN go build -o /opt/resource/out code.cloudfoundry.org/slack-attachment-resource/out
WORKDIR code.cloudfoundry.org/slack-attachment-resource

