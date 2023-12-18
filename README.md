# Snapping-thoughts (Twitter Bot cron-job)

Snapping-thoughts is a Twitter bot that tweets a random shower thought from the list of quotes (```src/seeds/item.seeds.go```) every 24 hours.

## Technologies

-   golang
-   AWS DynamoDB
-   AWS Lambda
-   Twitter API

## Getting Started

### Prerequisites

-   golang 1.21 or [later](https://go.dev)
-   makefile

### Installation

1. Clone this repository.
   ```bash
   git clone https://github.com/bookpanda/snapping-thoughts.git
   ```
2. Run `go mod download` to download all the dependencies.
3. Create a DynamoDB table with the following settings:
   - Table name: `snapping-thoughts`
   - Primary key: `Id` (String)
4. Copy `.env.template` and paste it in the same directory as `.env`, then fill in the appropriate values.
5. Make sure that in `main.go`, line 121 is commented out and line 117-118 is uncommented.
6. Run `make seed` to populate the table with the data.

### Running
-  Run `make tweet` or `go run src/main.go`

### Testing
1. Make sure that in `main.go`, line 121 is commented out and line 117-118 is uncommented.
2. Run `make test` or `go test  -v -coverpkg ./... -coverprofile coverage.out -covermode count ./...`

## Deployment
Consult [AWS Lambda deployment for Go](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html) to deploy the app to AWS Lambda. Make sure comment the correct line in main.go before deploying. Here are the commands to deploy to Lambda.
### Deploy with CloudFormation
1. Copy `example.deploy.yaml` and paste it in the same directory as `deploy.yaml`, then fill in the appropriate values. The ARN is the ARN of the IAM role for the Lambda function.
2. Run these:
```bash
make zip

# note that if you leave this in your bucket, you will be charged for the storage after some time
aws s3 cp myFunction.zip s3://<your-bucket-name>

aws cloudformation deploy \
  --template deploy.yaml \
  --stack-name snapping-thoughts
```

### Manual Deployment
1. Run these:
```bash
# for the first 2 commands, you can use 'make zip' instead
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go

zip myFunction.zip bootstrap

# the --role section requires your ARN of IAM role for that Lambda function. 
aws lambda create-function --function-name myFunction \
--runtime provided.al2023 --handler bootstrap \
--architectures arm64 \
--role arn:aws:iam::111122223333:role/lambda-ex \ \
--zip-file fileb://myFunction.zip

# you can use this to update the function code in Lambda
aws lambda update-function-code --function-name myFunction \
--zip-file fileb://myFunction.zip
```
2. After this, set the environment variables in the Lambda function. The environment variables are to be the same as the ones in your ```.env```.

After deployment, you can attach EventBridge to the Lambda function to run on a schedule. See [this](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-run-lambda-schedule.html) for more information.