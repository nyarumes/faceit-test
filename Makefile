SERVICE_NAME=service

run: build
	cd cmd/${SERVICE_NAME}; ./${SERVICE_NAME}

run/docker:
	

build:
	cd cmd/${SERVICE_NAME}; go build -o ${SERVICE_NAME} main.go

tests:
	go test ./...

tests/coverage:
	go test -cover ./...