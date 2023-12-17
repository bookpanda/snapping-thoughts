# Snapping-thoughts (Twitter Bot cron-job)

MyGraderList is a web app that lets students assess the difficulties and worthiness of each DSA grader problem in their respective courses.

MyGraderList Backend handles the business logic of the MyGraderList app i.e. CRUD operations for the problems' ratings, likes and emojis. 

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
1. Run `make test` or `go test  -v -coverpkg ./... -coverprofile coverage.out -covermode count ./...`