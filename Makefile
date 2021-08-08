.PHONY:
.SILENT:

build:
	go build -o ./main ./...

run: build
	./main

docker-build:
	docker build -t websocket-chat .

docker-run: docker-build
	docker run --name=chat -e CHAT_HOST=0.0.0.0 -e CHAT_PORT=8080 -p 8080:8080 --rm websocket-chat
