package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mailru/easyjson"
	"gopkg.in/yaml.v2"
	"log"
	"messagehook"
)

var (
	config *messagehook.Config
	blacklist *messagehook.Blacklist
)
func init() {
	bytes, err := Asset("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	patterns, err := config.LoadPatterns()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	blacklist = messagehook.NewBlacklist(patterns)
	log.Printf("blacklist is ready after loading %d patterns", len(patterns))
}

type Handler struct {}

func (h *Handler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	hookData := messagehook.Payload{}
	err := easyjson.Unmarshal(payload, &hookData)
	if err != nil {
		return nil, err
	}

	response := messagehook.Response{
		Message: hookData.Message,
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
