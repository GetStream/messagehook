package messagehook

import "github.com/GetStream/stream-chat-go/v2"

type Payload struct {
	Message stream_chat.Message `json:"message"`
	Channel stream_chat.Channel `json:"channel"`
}

type Response struct {
	Message stream_chat.Message `json:"message"`
}
