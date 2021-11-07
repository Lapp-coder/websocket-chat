package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"strings"

	"github.com/Lapp-coder/websocket-chat/internal/jrpc"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	url = "ws://localhost:8080/chat"
)

func main() {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	defer ws.Close()

	_, msg, err := ws.ReadMessage()
	if err != nil {
		logrus.Fatal(err)
	}

	client := jsonrpc.NewClient(ws.UnderlyingConn())
	defer client.Close()

	id := string(msg)
	go getUnreadMessages(id, client)

	fmt.Println(
		`
			 Введите id того пользоваетеля, которому хотите отправить сообщение. ( чтобы узнать свой id, введите команду /getmyid )
			 Вы также можете отправить сообщение сразу нескольким пользователям, просто укажите их id через запятую.
			 Для отправки сообщения всем пользователем введите *.
			 Для отправки сообщения самому себе введите echo.

			 Для того, чтобы отправить сообщение, выберете получателя (по принципу выше) и через %$4^b введите ваше сообщение.
			 Пример: echo %$4^b Hello, World!
	`)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		text := scanner.Text()

		switch text {
		case "/getmyid":
			fmt.Println("Ваш id:", id)
		default:
			sending := strings.Split(text, " %$4^b ")
			if len(sending) == 0 || len(sending) > 2 {
				logrus.Error("Неккоректный формат входных данных")
				continue
			}

			args := jrpc.SendMessageArgs{ID: id, IDs: sending[0], Message: sending[1]}
			if err = client.Call("Handler.SendMessage", args, nil); err != nil {
				logrus.Fatal(err)
			}
		}
	}
}

func getUnreadMessages(id string, client *rpc.Client) {
	for {
		args := jrpc.GetMessagesArgs{ID: id}
		var result *[]string
		if err := client.Call("Handler.GetMessages", args, &result); err != nil {
			logrus.Fatal(err)
		}

		if len(*result) != 0 {
			for _, msg := range *result {
				logrus.Println(msg)
			}
		}
	}
}
