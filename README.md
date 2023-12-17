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

2. Copy `.env.template` and paste it in the same directory as `.env`, then fill in the values.
3. Run `go mod download` to download all the dependencies.

### Running
-  Run `make tweet` or `go run src/main.go`

### Testing
1. Make sure that in ```main.go```, line 121 is commented out and line 118 is uncommented.
2. Run `make test` or `go test  -v -coverpkg ./... -coverprofile coverage.out -covermode count ./...`

### Deployment
Consult this [AWS Lambda deployment for Go](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html) to deploy the app to AWS Lambda. Make sure comment the correct line in main.go before deploying.