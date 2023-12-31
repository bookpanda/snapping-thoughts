tweet:
	go run main.go

seed:
	go run main.go seed

test:
	go vet ./...
	go test  -v -coverpkg ./src/... -coverprofile coverage.out -covermode count ./src/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

zip:
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
	zip myFunction.zip bootstrap