package messagehook

import "github.com/GetStream/stream-chat-go/v2"

type Config struct {
	MessageErrorText string `yaml:"message_error_text"`
	BlacklistRegexes []string `yaml:"blacklist_regexes"`
	MessageErrorAttachments []stream_chat.Attachment `yaml:"message_error_attachments"`
}
