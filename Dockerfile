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
RUN apk add --no-cache ca-certificates
RUN apk add mysql-client
COPY --from=builder /qiitaTrend ./
ENTRYPOINT ["./qiitaTrend"]
