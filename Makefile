build:
	go mod tidy
	go build -o scanner main.go
	go build -o keyid main_key_id.go
