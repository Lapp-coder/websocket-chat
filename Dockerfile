FROM golang:1.16.4-alpine3.13 AS builder

COPY ./ /github.com/Lapp-coder/websocket-chat
WORKDIR /github.com/Lapp-coder/websocket-chat

RUN go mod download
RUN go build -o ./build/.bin/chat ./cmd/app/server/main.go

FROM alpine:latest

WORKDIR /opt/

COPY --from=builder /github.com/Lapp-coder/websocket-chat/build/.bin/chat .

EXPOSE 8080

CMD ["./chat"]
