FROM golang:1.22-bullseye as builder
WORKDIR ${GOPATH}/src/github.com/arxon31/sso
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
      go build -ldflags "-w -s -X main.Build=`date +%FT%T%z`" cmd/sso/sso.go


FROM gcr.io/distroless/static-debian12 as prod
EXPOSE 8081
COPY --from=builder /go/src/github.com/arxon31/sso/sso /sso
ENTRYPOINT ["./sso"]
