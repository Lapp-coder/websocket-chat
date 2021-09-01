# Чат
[![Release](https://img.shields.io/badge/release-v1.1.0-blue)](https://github.com/Lapp-coder/websocket-chat/releases)

### При создании этого приложения были использованны следующие технологии:
* Go 1.16.6
* Протокол Websocket
* Протокол JSON-RPC
* Docker
* Git
* Makefile

### Для запуска сервера используйте следующую команду
```make docker-build && docker run --name=chat -e CHAT_HOST=<host> -e CHAT_PORT=<port> -p <port>>:<port> --rm websocket-chat```

### Для запуска сервера с конфигурацией по умолчанию
```make docker-run```

#### P.S: Для запуска клиента используйте команду ```make run-client```
