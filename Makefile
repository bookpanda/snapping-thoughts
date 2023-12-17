tweet:
	go run ./src/.

seed:
	go run ./src/. seed

test:
	go vet ./...
	go test  -v -coverpkg ./src/... -coverprofile coverage.out -covermode count ./src/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html