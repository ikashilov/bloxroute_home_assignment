.PHONY: build_server
build_server:
	cd server && \
	go build -o bin/server ./cmd/* && \
	cd ..

.PHONY: build_client
build_client:
	cd client && \
	go build -o bin/client ./cmd/* && \
	cd ..

.PHONY: build_all
build_all:
	make build_server & make build_client
