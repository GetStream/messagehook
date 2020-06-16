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

### Create Lambda Function

Let's asssume you want to call this function `message-hook`.

```
aws lambda create-function --function-name message-hook --runtime go1.x \
  --zip-file fileb://function.zip --handler main \
  --role arn:aws:iam::1234567:role/lambda-exec
```

### Update Lambda Function

```
aws lambda update-function-code --function-name message-hook --zip-file fileb://function.zip
```

### Test Lambda Execution

You can test the Lambda using aws lambda from awscli.

```
aws lambda invoke --function-name message-hook --payload '{ "message": {"text": "example@gmail.com"} }' response.json; cat response.json | jq
```