### Multi-stage build
FROM jormungandrk/goa-build as build

COPY . /go/src/github.com/JormungandrK/jwt-issuer
RUN go install github.com/JormungandrK/jwt-issuer


### Main
FROM alpine:3.6

COPY --from=build /go/bin/jwt-issuer /usr/local/bin/jwt-issuer
EXPOSE 8080

ENV API_GATEWAY_URL="http://localhost:8001"

CMD ["/usr/local/bin/jwt-issuer"]
