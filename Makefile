build:
	echo "Building Package"
	go build -o bin/s3analytics cmd/main.go

test:
	echo "Testing all packages with cover"
	go test ./... -cover

run:
	echo "Running s3analytics"
	go run cmd/main.go

compile: 
	echo "Compiling for all required OS and Platform"
	GOOS=linux GOARCH=amd64 go build -o bin/s3analytics-linux-amd64 cmd/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/s3analytics-windows-amd64.exe cmd/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/s3analytics-darwin-amd64 cmd/main.go
	
all: build test compile