BINARY_CEP=bin/cep

.PHONY: cep
cep: clean tidy
	GOOS=linux go build -o ./$(BINARY_CEP)
	./$(BINARY_CEP)

.PHONY: run
run: clean tidy
	go run -race main.go

clean:
	rm -rf bin/*

.PHONY: test
test:
	go test ./... -race -cover

.PHONY: tidy 
tidy:
	go mod tidy
