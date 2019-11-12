all: build test

fetch_dependencies:
	go get -v -t -d ./...

build: fetch_dependencies
	go build -v -o .bin/container-linux-node-labeler .

test: fetch_dependencies
	go test -v ./...

run: build
	./.bin/container-linux-node-labeler --kubeconfig ~/.kube/config

docker_image:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o .bin/container-linux-node-labeler-docker .
	docker build . -t container-linux-node-labeler:0.0.1
