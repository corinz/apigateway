clean:
	go clean
	go fmt ./...
	rm -f .cert/localhost*
	openssl req -x509 -newkey rsa:4096 -sha256 -nodes -keyout .cert/localhost.key -out .cert/localhost.crt -days 365 -subj '/CN=localhost'

run:
	go fmt ./...
	go run cmd/server/main.go  
