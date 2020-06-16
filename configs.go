package messagehook

import (
	"bufio"
	"log"

	stream_chat "github.com/GetStream/stream-chat-go/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"
)

type Config struct {
	MessageErrorText        string                   `yaml:"message_error_text"`
	Patterns                []string                 `yaml:"patterns"`
	S3Bucket                string                   `yaml:"s3_bucket"`
	S3File                  string                   `yaml:"s3_file"`
	S3Region                string                   `yaml:"s3_region"`
	MessageErrorAttachments []stream_chat.Attachment `yaml:"message_error_attachments"`
	StreamApiKey            string                   `yaml:"stream_api_key"`
	StreamApiSecret         string                   `yaml:"stream_api_secret"`
	StreamBaseURL           string                   `yaml:"stream_base_url"`
	CheckSignature          bool                     `yaml:"check_signature"`
}

func NewFromBytes(bytes []byte) (*Config, error) {
	config := Config{}
	err := yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	patterns, err := config.LoadPatterns()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	config.Patterns = patterns
	return &config, nil
}

func LoadFromS3(bucket, path, region string) ([]string, error) {
	s, err := session.NewSession(
		&aws.Config{Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}
	svc := s3.New(s)
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	}
	out, err := svc.GetObject(input)
	if err != nil {
		return nil, err
	}

	defer out.Body.Close()

	var lines []string
	scanner := bufio.NewScanner(out.Body)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func (c Config) LoadPatterns() ([]string, error) {
	if c.S3Bucket != "" && c.S3File != "" && c.S3Region != "" {
		return LoadFromS3(c.S3Bucket, c.S3File, c.S3Region)
	}

	return c.Patterns, nil
}
