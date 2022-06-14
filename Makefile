SERVICE_NAME=service

run: build
	cd cmd/${SERVICE_NAME}; ./${SERVICE_NAME}

build:
	cd cmd/${SERVICE_NAME}; go build -o ${SERVICE_NAME} main.go

tests:
	go test ./...

tests/integration:


tests/coverage:
	go test -cover ./...