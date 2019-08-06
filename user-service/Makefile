include .env

.PHONY: run_test
run_test:
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "test report stored in coverage.html"


.PHONY: build
build:
	@echo "Building binary..."
	go build -o app

.PHONY: docker_build
docker_build: build
	@echo "Building docker image..."
	docker build -t shuzasa/user-service:v1.0 .

.PHONY: docker_run
docker_run: docker_build
	docker run -p 8080:8080 --env-file .env shuzasa/user-service:v1.0