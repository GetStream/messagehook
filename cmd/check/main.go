package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"messagehook"
	"os"

	"github.com/mailru/easyjson"
)

func main() {
	reader, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	config, err := messagehook.NewFromBytes(bytes)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	blacklist := messagehook.NewBlacklist(config.Patterns)
	log.Printf("blacklist is ready after loading %d patterns", len(config.Patterns))

	stdinBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	hookData := messagehook.Payload{}

	err = easyjson.Unmarshal(stdinBytes, &hookData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	response := messagehook.Response{
		Message: hookData.Message,
	}

	if blacklist.Match(response.Message.Text) {
		messagehook.RewriteMessageAsError(&response.Message, config.MessageErrorText, config.IncludeOriginalText, config.MessageErrorAttachments)
	}

	out, err := easyjson.Marshal(response)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println(string(out))
}
