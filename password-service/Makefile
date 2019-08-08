include .env

.PHONY: run_test
run_test:
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "test report stored in coverage.html"

.PHONY: run
run:
	@go run main.go

.PHONY: build
build:
	@echo "Building binary..."
	GOOS=linux go build -o app

.PHONY: docker_build
docker_build: build
	@echo "Building docker image..."
	docker build -t shuzasa/password-service:v1.0 .

.PHONY: docker_push
docker_push: docker_build
	@echo "Pushing docker image..."
	docker push shuzasa/password-service:v1.0

.PHONY: docker_run
docker_run: docker_build
	docker run -p 8081:8081 --env-file .env shuzasa/password-service:v1.0

.PHONY: deploy
deploy:
	kubectl apply -f ./deployment/password.yaml

.PHONY: delete
delete:
	kubectl delete svc password-service
	kubectl delete deployment password-deploy