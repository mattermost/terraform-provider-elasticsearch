TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=mattermost.com
NAMESPACE=terraform
NAME=elasticsearch
BINARY=terraform-provider-${NAME}
VERSION=0.1
OS_ARCH=darwin_amd64

default: install # Installs the provider in .terraform.d/plugins

build: # Builds the binary for local testing
	@echo Build binary for local
	go build -o ${BINARY}

release: # Builds the binary for provider
	@echo Building Provider
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

run-elasticsearch: # Run the docker-compose to test
	@echo Running elasticsearch
	docker-compose -f docker/docker-compose.yml up -d
	sleep 15

install: build # Builds and creates the necessary dirs to add the plugin
	@echo Installing provider in Terraform
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

install-deps: # Install golangci-lint
	@echo Installing Golang-CI Linter
	./scripts/install-ci.sh

lint: # Run lint
	@echo Run linter
	golangci-lint run --timeout=5m

test: # Run unit tests
	@echo Run tests
	go test $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4               

testacc: # Run acceptance tests
	@echo Run acceptance tests
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
