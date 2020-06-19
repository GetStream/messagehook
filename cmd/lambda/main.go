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
	config      *messagehook.Config
	blacklist   *messagehook.Blacklist
	chatClients map[string]*stream_chat.Client
)

const (
	apiSecretHeadeName = "x-signature"
	apiKeyHeaderName   = "x-api-key"
)

func init() {
	chatClients = map[string]*stream_chat.Client{}

	bytes, err := Asset("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	config, err = messagehook.NewFromBytes(bytes)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, credentials := range config.StreamCredentials {
		log.Printf("loading credentials for api key: %q", credentials.Key)
		chatClient, err := stream_chat.NewClient(credentials.Key, []byte(credentials.Secret))
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		if config.StreamBaseURL != "" {
			chatClient.BaseURL = config.StreamBaseURL
		}
		chatClients[credentials.Key] = chatClient
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

func parseGatewayPayload(data []byte) (*messagehook.Payload, error) {
	var body []byte
	event := apiGatewayEvent{}

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

	payload := messagehook.Payload{}

	err = easyjson.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	if config.CheckSignature {
		key := event.Headers[apiKeyHeaderName]
		secret := event.Headers[apiSecretHeadeName]

		client, found := chatClients[key]
		if !found {
			return nil, fmt.Errorf("handler is not configured for api key %q", key)
		}

		if !client.VerifyWebhook(body, []byte(secret)) {
			return nil, fmt.Errorf("payload is not signed with the correct secret for api key %q", secret)
		}
	}

	return &payload, nil
}

func (h *Handler) Invoke(_ context.Context, payload []byte) ([]byte, error) {
	request, err := parseGatewayPayload(payload)
	if err != nil {
		return nil, err
	}

	response := messagehook.Response{
		Message: request.Message,
	}

	if blacklist.Match(response.Message.Text) {
		log.Printf("message %q matched the blacklist and will be rewritten", response.Message.Text)
		messagehook.RewriteMessageAsError(&response.Message, config.MessageErrorText, config.IncludeOriginalText, config.MessageErrorAttachments)
	}

	return easyjson.Marshal(response)
}

func main() {
	handler := &Handler{}
	lambda.StartHandler(handler)
}
