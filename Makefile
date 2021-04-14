all: clean build docker
clean:
	go clean
	go fmt ./...
	rm -f .cert/localhost*
	openssl req -x509 -newkey rsa:4096 -sha256 -nodes -keyout .cert/localhost.key -out .cert/localhost.crt -days 365 -subj '/CN=localhost'

run:
	go fmt ./...
	go run cmd/server/main.go  

build:
	go fmt ./...
	CGO_ENABLED=0 GOOS=darwin go build -o bin/apigateway-darwin cmd/server/main.go
	CGO_ENABLED=0 GOOS=linux go build -o bin/apigateway cmd/server/main.go
	CGO_ENABLED=0 GOOS=windows go build -o bin/apigateway.exe cmd/server/main.go

docker:
	go fmt ./...
	@[ "${REG}" ] || read -p "Container tag {repo}/{image}:{tag}: " REG \
	&& echo $$REG \
	&& docker build -t $$REG . \
	&& docker push $$REG
