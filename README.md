# go-openapi ✨

This project shows a basic Golang server with a Dockerfile.
It uses go-swagger to generate the code from the swagger file which is cool. 😎


It can be consumed and automatically deployed to Kubernetes ( with generation of charts ) through gitlab.

- OpenAPI
- Kubernetes
- Docker
- Gitlab
- go-memdb ( in memory db )
- AutoDevops

### Regenerate the code 💅🏼

This project uses OPENAPI2.0 this means you can generate all the stub code from the spec.

```
brew tap go-swagger/go-swagger
brew install go-swagger

swagger generate server -f static/swagger.yaml -A go-openapi --exclude-main
```

Running...

```
go run cmd/go-openapi/main.go
```

View the UI...

```
http://127.0.0.1:58845/v2/docs
```

## Do something 🤷🏼‍♀️
```bash
# Create a user
 curl -X POST "http://localhost:8080/v2/user" -H  "accept: application/xml" -H  "Content-Type: application/json" -d "{  \"email\": \"string\",  \"firstName\": \"string\",  \"id\": 0,  \"lastName\": \"string\",  \"password\": \"string\",  \"phone\": \"string\",  \"userStatus\": 0,  \"username\": \"alex\"}"
# Get that user
curl -X GET "http://localhost:8080/v2/user/alex" -H  "accept: application/xml" -v
```


#### Run Jaeger 

Set `JAEGER_AGENT_HOST` and `JAEGER_AGENT_PORT` otherwise you can run locally with docker...

```
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.22
```