all: build test

fetch_dependencies:
	go get -v -t -d ./...

build: fetch_dependencies
	go build -v -o .bin/container-linux-node-labeller .

test: fetch_dependencies
	go test -v ./...

run: build
	./.bin/container-linux-node-labeller --kubeconfig ~/.kube/config

docker_image:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o .bin/container-linux-node-labeller-docker .
	docker build . -t mwarzynski/container-linux-node-labeller:0.0.1

docker_push: docker_image
	docker push mwarzynski/container-linux-node-labeller:0.0.1
