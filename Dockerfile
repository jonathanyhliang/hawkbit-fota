FROM golang:alpine

COPY ../ /root/hawkbit-fota

WORKDIR /root/hawkbit-fota

RUN go mod tidy \
    && go build .
