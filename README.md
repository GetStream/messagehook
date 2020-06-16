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

### Update Lambda Function

```
aws lambda update-function-code --function-name message-hook --zip-file fileb://function.zip
```

### Test Lambda Execution

You can test the Lambda using aws lambda from awscli.

```
aws lambda invoke --function-name message-hook --payload '{ "message": {"text": "example@gmail.com"} }' response.json; cat response.json | jq
```

### Check locally

Testing Lambdas can be annoying especially if you do not have a working setup for serverless. You can use the check command to do some basic checks.

```
echo '{"message": {"text": "tbarbugli@gmail.com"}}' | go run cmd/check/main.go
```
