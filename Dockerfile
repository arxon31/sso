FROM golang:1.22-bullseye as builder
WORKDIR ${GOPATH}/src/github.com/arxon31/sso
COPY . .
RUN go build -ldflags "-w -s -X main.Build=`date +%FT%T%z`" cmd/sso/sso.go


FROM debian:bullseye-slim as prod
EXPOSE 8081
COPY --from=builder /go/src/github.com/arxon31/sso/sso /sso
ENTRYPOINT ["/sso"]
