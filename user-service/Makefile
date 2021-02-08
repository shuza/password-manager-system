include .env

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
COVERAGE_MIN=80

.PHONY: test
test:
	@echo "> running test and creating coverage report"
	go test -race -p=1 -cover -coverprofile=coverage.out $(ALL_PACKAGES)
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | grep -i total:
	@go tool cover -func=coverage.out | gawk '/total:.*statements/ {if (strtonum($$3) < $(COVERAGE_MIN)) {print "ERR: coverage is lower than $(COVERAGE_MIN)"; exit 1}}'
	@rm *.out

.PHONY: run
run:
	@go run main.go server

.PHONY: build
build:
	@echo "Building binary..."
	GOOS=linux go build -o app

.PHONY: docker_build
docker_build: build
	@echo "Building docker image..."
	docker build -t shuzasa/user-service:v1.0 .

.PHONY: docker_push
docker_push: docker_build
	@echo "Pushing docker image..."
	docker push shuzasa/user-service:v1.0

.PHONY: docker_run
docker_run: docker_build
	docker run -p 8080:8080 --env-file .env shuzasa/user-service:v1.0

.PHONY: deploy
deploy:
	kubectl apply -f ./deployment/user.yaml

.PHONY: delete
delete:
	kubectl delete svc user-service
	kubectl delete deployment user-deploy
