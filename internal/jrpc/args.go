package jrpc

type SendMessageArgs struct {
	ID      string
	IDs     string
	Message string
}

type GetMessagesArgs struct {
	ID string
}
