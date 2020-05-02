.PHONY: run
run:
	docker network create user-api-processor || true
	docker-compose up -d --build

.PHONY: stop
stop:
	docker-compose down