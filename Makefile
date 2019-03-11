## Makefile

build: ## Builds the starter pack
	go build -o bin/broker ./cmd/broker

build-linux: ## Builds a Linux executable
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -o bin/broker-linux --ldflags="-s" ./cmd/broker

build-docker: ## Builds a docker image
	@docker build . -t alekssaul/sb-router:dev