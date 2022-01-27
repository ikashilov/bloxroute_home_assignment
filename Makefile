.PHONY: build_server
build_server:
	go build -o bin/server ./cmd/server/*

.PHONY: build_client
build_client:
	go build -o bin/client ./cmd/client/*

.PHONY: build_all
build_all:
	make build_server & make build_client
