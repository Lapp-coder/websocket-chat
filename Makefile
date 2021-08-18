.PHONY:
.SILENT:

build-server:
	go build -o ./build/.bin/server ./cmd/app/server/main.go

run-server: build-server
	./build/.bin/server

build-client:
	go build -o ./build/.bin/client ./cmd/app/client/main.go

run-client: build-client
	./build/.bin/client

docker-build:
	docker build -t websocket-chat .

docker-run: docker-build
	docker run --name=chat -e CHAT_HOST=0.0.0.0 -e CHAT_PORT=8080 -p 8080:8080 --rm websocket-chat
