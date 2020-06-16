package messagehook

import "github.com/GetStream/stream-chat-go/v2"

func RewriteMessageAsError(message *stream_chat.Message, errorMessage string) {
	message.Text = errorMessage
	message.Type = stream_chat.MessageTypeError
}
