all: clean build docker swagger
clean:
	go clean
	go mod tidy
	go fmt ./...
	rm -f cmd/server/apigateway.json

cert-renew:
	rm -f .cert/localhost*
	openssl req -x509 -newkey rsa:4096 -sha256 -nodes -keyout .cert/localhost.key -out .cert/localhost.crt -days 365 -subj '/CN=localhost'

run:
	go fmt ./...
	go run cmd/server/main.go  

build: clean
	CGO_ENABLED=0 GOOS=darwin go build -o bin/apigateway-darwin cmd/server/main.go
	CGO_ENABLED=0 GOOS=linux go build -o bin/apigateway cmd/server/main.go
	CGO_ENABLED=0 GOOS=windows go build -o bin/apigateway.exe cmd/server/main.go

swagger-check:
	which swagger || (go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger-gen:
	swagger generate spec -o ./swagger.json --scan-models

swagger: swagger-check swagger-gen
	swagger validate swagger.json && mv swagger.json docs/swagger-ui/dist/swagger.json

swagger-serve: swagger-check swagger-gen
	swagger validate swagger.json && swagger serve -F=swagger swagger.json

docker: swagger clean
	@[ "${REG}" ] || read -p "Container tag {repo}/{image}:{tag}: " REG \
	&& echo $$REG \
	&& docker build -t $$REG . \
	&& docker push $$REG
