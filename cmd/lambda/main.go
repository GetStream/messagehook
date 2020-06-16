package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"messagehook"

	stream_chat "github.com/GetStream/stream-chat-go/v2"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mailru/easyjson"
)

var (
	config     *messagehook.Config
	blacklist  *messagehook.Blacklist
	chatClient *stream_chat.Client
)

const (
	apiSecret = "x-signature"
)

func init() {
	bytes, err := Asset("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	config, err = messagehook.NewFromBytes(bytes)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("api secret: %q", config.StreamApiSecret)

	chatClient, err = stream_chat.NewClient(config.StreamApiKey, []byte(config.StreamApiSecret))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if config.StreamBaseURL != "" {
		chatClient.BaseURL = config.StreamBaseURL
	}

	blacklist = messagehook.NewBlacklist(config.Patterns)
	log.Printf("blacklist is ready after loading %d patterns", len(config.Patterns))
}

type Handler struct{}

type apiGatewayEvent struct {
	Headers map[string]string `json:"headers"`
	Base64  bool              `json:"isBase64Encoded"`
	Body    string            `json:"body"`
}

type testSignature func(body, signature []byte) (valid bool)

func noopCheck(_, _ []byte) bool { return true }

func parseGatewayPayload(data []byte, check testSignature) (*messagehook.Payload, error) {
	var body []byte
	event := apiGatewayEvent{}

	log.Printf("request body %s", string(data))

	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	body = []byte(event.Body)

	if event.Base64 {
		body, err = base64.StdEncoding.DecodeString(event.Body)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("%q", string(body))
	log.Printf("%q", event.Headers[apiSecret])

	payload := messagehook.Payload{}

	err = easyjson.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	if !check(body, []byte(event.Headers[apiSecret])) {
		return nil, fmt.Errorf("payload is not signed with the same secret")
	}

	return &payload, nil
}

func (h *Handler) Invoke(_ context.Context, payload []byte) ([]byte, error) {
	check := noopCheck
	if config.CheckSignature {
		check = chatClient.VerifyWebhook
	}

	request, err := parseGatewayPayload(payload, check)
	if err != nil {
		return nil, err
	}

	response := messagehook.Response{
		Message: request.Message,
	}

	if blacklist.Match(response.Message.Text) {
		messagehook.RewriteMessageAsError(&response.Message, config.MessageErrorText)
	}

	return easyjson.Marshal(response)
}

func main() {
	handler := &Handler{}
	lambda.StartHandler(handler)
}
