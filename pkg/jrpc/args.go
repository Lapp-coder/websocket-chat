package jrpc

type SendMessageArgs struct {
	ID      string
	Ids     string
	Message string
}

type GetMessagesArgs struct {
	ID string
}
