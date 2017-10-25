### Multi-stage build
FROM golang:1.8.3-alpine3.6 as build

RUN apk --no-cache add git curl openssh

COPY keys/id_rsa /root/.ssh/id_rsa
RUN chmod 700 /root/.ssh/id_rsa && \
    echo -e "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config && \
    git config --global url."ssh://git@github.com:".insteadOf "https://github.com"

RUN go get -u github.com/goadesign/goa/... && \
    go get -u gopkg.in/mgo.v2 && \
    go get -u golang.org/x/crypto/bcrypt && \
    go get -u github.com/afex/hystrix-go/hystrix && \
    go get -u github.com/satori/go.uuid && \
    go get -u github.com/dgrijalva/jwt-go

RUN go get -u github.com/JormungandrK/microservice-tools; \
    go get -u github.com/JormungandrK/microservice-security/...; \
    exit 0

COPY . /go/src/github.com/JormungandrK/jwt-issuer
RUN go install github.com/JormungandrK/jwt-issuer


### Main
FROM alpine:3.6

COPY --from=build /go/bin/jwt-issuer /usr/local/bin/jwt-issuer
EXPOSE 8080

ENV API_GATEWAY_URL="http://localhost:8001"

CMD ["/usr/local/bin/jwt-issuer"]
