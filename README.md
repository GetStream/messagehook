### Message Hook Handler

This is an example Lambda handler that you can use as a [before message send](https://getstream.io/chat/docs/) hook endpoint. The handler will receive a message
payload and checks if the message text matches any regex. If there is a match the message will be rewritten as an error message.

This handler also checks that the message includes a correct signature, the signature can only be created using the API secret. 

You can configure several things to match your use-case but it is very trivial to perform different logic by making minimal changes to the code (see cmd/lambda/main.go)

### Install

Make sure you have the latest version of Go installed, once you have that you only need to run `make` and you will get a 
ready to use lambda function (function.zip).

Note: `make` will fail if you `config.yaml` file does not exist.

You also need to have `awscli` installed on your laptop.

### Configuration

Before using the lambda you need to create your own config YML file. This configuration file contains the list of
regexes your hook should match as well as the error message to use for the rewritten message.

```bash
cp example.config.yaml config.yaml
```

The configuration file is embedded in the final lambda function. Make sure to run `make` every time you change the configuration file.

### API Credentials & Signature check

Make sure you use the correct API Keys and to enable signature checking.

```yaml
stream_api_key: api-key
stream_api_secret: api-secret
check_signature: true
```

### Configuring Patterns

The list of patterns to match messages can be provided in two ways: via the configuration file or by storing it on S3.

#### S3 configuration

When the S3 configuration is provided, the lambda function will read the S3 bucket and load all patterns in the file

```yaml
s3_bucket: my-bucket-name
s3_file: path/to/file
s3_region: aws-region-name
```

If you decide to use this approach, make sure that the Lambda function role can perform GetObject to your S3 bucket.

#### YAML

If S3 is not used, the lambda will load the list of patterns from the YAML file itself

```yaml
patterns:
  - "[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*" # emails
```

### Create Lambda Function

Let's asssume you want to call this function `message-hook`.

```
aws lambda create-function --function-name message-hook --runtime go1.x \
  --zip-file fileb://function.zip --handler main \
  --role arn:aws:iam::1234567:role/lambda-exec
```

### Lambda Invoke

The lambda function itself does not work if invoked using the Lambda API and it is meant to be used as the handler for an API gateway route.
Using something like `aws lambda invoke` will return an error since the payload is not the same as what API Gateway fordwards.

### Update Lambda Function

```
aws lambda update-function-code --function-name message-hook --zip-file fileb://function.zip
```

### Check locally

Testing Lambdas can be annoying especially if you do not have a working setup for serverless. You can use the check command to do some basic checks.

```
echo '{"message": {"text": "tbarbugli@gmail.com"}}' | go run cmd/check/main.go
```

### API Gateway

Once you have the lambda function up and running, make sure to create an API Gateway and to connect the lambda with a route.

1. Go to API Gateway https://console.aws.amazon.com/apigateway
1. Create a new API
1. Select the HTTP API type
1. Choose Lambda from the list of integrations and select the Lambda function created earlier
1. Make sure to select version 2.0 or the Lambda will not work correctly 
![image](https://user-images.githubusercontent.com/88735/84822134-824eda00-b01c-11ea-95af-63c9502b4532.png)
1. Create a new route that sends POST requests to the Lambda function
![image](https://user-images.githubusercontent.com/88735/84822608-3d777300-b01d-11ea-9a47-871a95c04ca6.png)
1. Continue with the process until the API is created
1. The URL of the API gateway will be visible at the end
1. Make sure to build the full URL (Invoke URL + Route path) ie. `	https://fuhlrrq4h0.execute-api.us-east-1.amazonaws.com/message-hook`

### Configure the hook on your Stream Application

```js
async function main() {
    const client = new StreamChat(apiKey, apiSecret);
    await client.updateAppSettings({
        before_message_send_hook_url: 'http://127.0.0.1:4323',
    });
}
```