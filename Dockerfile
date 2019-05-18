### Multi-stage build
FROM golang:1.10-alpine3.7 as build

RUN apk --no-cache add git curl openssh

RUN go get -u -v github.com/keitaroinc/goa/... && \
    go get -u -v github.com/afex/hystrix-go/hystrix && \
    go get -u -v github.com/Microkubes/microservice-security/... && \
    go get -u -v github.com/Microkubes/microservice-tools/...

COPY . /go/src/github.com/Microkubes/jwt-issuer

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install github.com/Microkubes/jwt-issuer


### Main
FROM scratch

ENV API_GATEWAY_URL="http://localhost:8001"

COPY --from=build /go/src/github.com/Microkubes/jwt-issuer/config.json /config.json
COPY --from=build /go/bin/jwt-issuer /jwt-issuer

EXPOSE 8080

CMD ["/jwt-issuer"]
