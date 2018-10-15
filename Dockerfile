FROM golang:latest as builder
WORKDIR /go/src/github.com/jskelcy/access-stats
COPY . .
RUN make build-linux

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/jskelcy/access-stats/out .
RUN touch /var/log/access.log
CMD ["./access-stats"]  