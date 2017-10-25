JWT Token Issuer
================

[![Build](https://travis-ci.com/JormungandrK/jwt-issuer.svg?token=UB5yzsLHNSbtjSYrGbWf&branch=master)](https://travis-ci.com/JormungandrK/jwt-issuer)
[![Test Coverage](https://api.codeclimate.com/v1/badges/7b8eb0b65c625e8ceb7c/test_coverage)](https://codeclimate.com/repos/59e7253fb82c7d02d200155a/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/7b8eb0b65c625e8ceb7c/maintainability)](https://codeclimate.com/repos/59e7253fb82c7d02d200155a/maintainability)

Microservice that exposes endpoint for issuing new JWT tokens.

#Issuing a JWT tokens

A client can get a JWT token by accesing the signin endpoint at ```<jwt-issuer-host:port>/jwt/signin```.
The client must send a POST request (content type ```application/x-www-form-urlencoded``` - form post)
with the following parameters:
 * ```username``` - the user's username
 * ```password``` - password
 * ```scope``` - the scope for the request (```api:read``` or ```api:write```)

An example with ```curl```:

```bash
curl -v -X POST -d "username=user&password=p@ss&scope=api:read" "http://jwt.myhost:8080/jwt/signin"

> POST /jwt/signin HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.47.0
> Accept: */*
> Content-Length: 47
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 47 out of 47 bytes
< HTTP/1.1 201 Created
< Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDMwNTIwNDUsImlhdCI6MTUwMzA1MjAxNSwiaXNzIjoiSm9ybXVuZ2FuZHIgSldUIEF1dGhvcml0eSIsImp0aSI6ImQ4ZWU3NTRkLTc3YzAtNDBlOS1iN2ViLWRiY2Y1ZjVhMDlhZSIsIm5iZiI6MCwib3JnYW5pemF0aW9ucyI6IiIsInJvbGVzIjoidXNlciIsInNjb3BlcyI6ImFwaTpyZWFkIiwic3ViIjoiNTk5NDFjNWQwMDAwMDAwMDAwMDAwMDAwIiwidXNlcklkIjoiNTk5NDFjNWQwMDAwMDAwMDAwMDAwMDAwIiwidXNlcm5hbWUiOiJwYXZsZSJ9.HB7G5LXZgGK0wxLHIndtU_iwlzJP1ptDY2yhR7bADTB7kw0I8pU579QY5cr_tpc1GwTb3iev1pZvpB_XnNHRQonp6JIzeIUAFvZI4-X-fq7S_sfEMQyK12Id6sIr1MoIoFYPmgQGRlh5hJMWNS9UdeQp8qqAMQvEx42qCtrRUI_wQDl48V_Yp_fn_82DWWJZFEJ4FLfKu5l6bkJWpYcj3ChF-OrjP2uMcjMU1s3vUEnn6w9QuEgY1lYLjzMxVjDD0fTknNERrCaXFS25wbZl2WQYq62OcDsU1vjMCf_n3aPxP_He_I4nabJWtdIltoJC6UH-z5AZEUClFJs1sbYKEA
< Date: Fri, 18 Aug 2017 10:26:55 GMT
< Content-Length: 0
< Content-Type: text/plain; charset=utf-8
<

```


The JWT token will be available in the HTTP Response, the value of the ```Authorization``` header.

# Service Configuration

The service is configured using a JSON file with the following structure:
```javascript
{
  "jwt":{ // JWT Configuration
    "issuer": "Jormungandr JWT Authority", // The name of the JWT issuer
    "signingMethod": "RS512", // Method used for signing (RS256, RS512 etc)
    "expiryTime": 30000 // JWT token validity period. The token expires in this many milliseconds after its being generated.
  },
  "keys": { // Map of keys. Must contain at least "default" and "system".
    "default": "./test-keys/rsa_default", // Used for generating and signing the JWT tokens for the clients.
    "system": "./test-keys/rsa_system" // Used for JWT token for accesing the User Microservice internally.
  },
  "microservice": { // Microservice configuration
    "name": "jwt-issuer",
    "port": 8080,
    "virtual_host": "jwt.auth.jormugandr.org",
    "hosts": ["localhost", "jwt.auth.jormugandr.org"],
    "weight": 10,
    "slots": 100
  },
  "services": { // Map of URLs for the internal services
    "user-microservice": "http://kong.gateway:8001/user" // MUST contain URL for the "user-microservice". Set this to the Kogn API Gateway URL for the user microservice.
  }
}

```

Default path is /run/secrets/microservice_apps_management_config.json. To change the path set the **SERVICE_CONFIG_FILE** env var.
