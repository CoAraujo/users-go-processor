.PHONY: build-go-processor
build-go-processor:
	docker build -t go-processor .

.PHONY: up
up:
	docker-compose up -d

.PHONY: all
all:
	build-go-processor up
