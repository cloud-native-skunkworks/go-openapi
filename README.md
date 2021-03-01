# go-openapi

This project shows a basic Golang server with a Dockerfile.

It can be consumed and automatically deployed to Kubernetes ( with generation of charts ) through gitlab.

- OpenAPI
- Kubernetes
- Docker
- Gitlab
- AutoDevops

### Regenerate the code

This project uses OPENAPI2.0 this means you can generate all the stub code from the spec.

```
brew tap go-swagger/go-swagger
brew install go-swagger

swagger generate server -f static/swagger.yaml -A gitlab-auto-devops-example --exclude-main
```

Running...

```
go run cmd/gitlab-auto-devops-example-server/main.go
```

Test the health check...

```
curl http://127.0.0.1:8080d/v2/healthz -v
```
