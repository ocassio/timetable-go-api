ARG GOLANG_VERSION="1.14"

FROM golang:$GOLANG_VERSION-alpine as builder
WORKDIR /go/src/github.com/ocassio/timetable-go-api
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s' -o ./timetable-api

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/ocassio/timetable-go-api /
ENTRYPOINT ["/timetable-api"]
