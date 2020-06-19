package messagehook

import stream_chat "github.com/GetStream/stream-chat-go/v2"

const (
	originalTextField = "original_text"
	isRewritten       = "rewritten_by_presend_hook"
)

func RewriteMessageAsError(message *stream_chat.Message, errorMessage string, includeOriginal bool, attachments []stream_chat.Attachment) {
	if includeOriginal {
		if message.ExtraData == nil {
			message.ExtraData = map[string]interface{}{}
		}
		message.ExtraData[originalTextField] = message.Text
	}
	message.ExtraData[isRewritten] = true
	message.Text = errorMessage
	for i := range attachments {
		message.Attachments = append(message.Attachments, &attachments[i])
	}
}
