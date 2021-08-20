FROM golang:1.16-alpine as builder

WORKDIR $GOPATH/src/kts/push
COPY . .
RUN ls -al .
RUN go build pkg/main.go
RUN cp main /

FROM alpine:3.14

COPY --from=builder /main /app/
WORKDIR /app
CMD ./main
