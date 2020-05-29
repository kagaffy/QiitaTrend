FROM golang:alpine as builder
ENV APPDIR $GOPATH/src/github.com/kagaffy/QiitaTrend
ENV GO111MODULE on
RUN \
  apk update --no-cache && \
  mkdir -p $APPDIR
ADD . $APPDIR/
WORKDIR $APPDIR
RUN go build -ldflags "-s -w" -o qiitaTrend qiitaTrend.go
RUN mv qiitaTrend /

FROM alpine
#RUN apk add --no-cache ca-certificates
RUN apk add --no-cache udev ttf-freefont chromium chromium-chromedriver
RUN apk add mysql-client
EXPOSE 8080
COPY --from=builder /qiitaTrend ./
ENTRYPOINT ["./qiitaTrend"]

#FROM golang:1.14-alpine AS build-env
#
#ENV APPDIR $GOPATH/src/github.com/kagaffy/QiitaTrend
#WORKDIR $APPDIR
#COPY . .
#
#RUN go build
#
#FROM alpine:latest
##LABEL maintainer="your name <foo@example.com>"
#
#WORKDIR /project
#COPY --from=build-env /go/src/path/to/your/project .
#
#RUN apk add --no-cache udev ttf-freefont chromium chromium-chromedriver
#
#EXPOSE 8080
#
#ENTRYPOINT ["./qiitaTrend"]
