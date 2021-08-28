FROM golang:1.16-alpine as builder

WORKDIR $GOPATH/src/github.com/ktsstudio/selectel-exporter
COPY . .
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/selectel-exporter pkg/main.go

FROM scratch
COPY --from=builder /go/bin/selectel-exporter /go/bin/selectel-exporter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/selectel-exporter"]
