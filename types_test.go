package messagehook

import (
	"github.com/GetStream/stream-chat-go/v2"
	"github.com/mailru/easyjson"
	"testing"
)

func TestMessage_Marshal(t *testing.T) {
	message := stream_chat.Message{
		Text: "heloooo",
		ExtraData: map[string]interface{}{
			"custom": true,
		},
	}
	_, err := easyjson.Marshal(message)
	if err != nil {
		t.Error(err)
	}
}
