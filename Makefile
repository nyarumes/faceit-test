USER_SERVICE=user

run: build
	cd cmd/${USER_SERVICE}; ./${USER_SERVICE}

build:
	cd cmd/${USER_SERVICE}; go build -o ${USER_SERVICE} main.go