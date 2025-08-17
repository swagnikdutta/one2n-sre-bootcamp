include .env

run:
# just cmd does not work as go interprets it to be package path, containing multiple source files/sub-packages even
# if it doesn't. So we add ./ so that it is treated as an import path.
	go run ./cmd
build:
	go build -o main ./cmd
fmt:
	go fmt ./...
vet:
	go vet ./...
clean:
	rm -f main
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
test:
	go test ./...

migrate-up:
	migrate -database "$(DATABASE_URL)" -path migrations up
migrate-down:
	migrate -database "$(DATABASE_URL)" -path migrations down
db-up:
	docker compose up -d db
db-down:
	docker compose down db
docker-build-image:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f alpine.Dockerfile .
docker-run-backend:
	# commenting the below command and using docker compose instead as only compose can ensure that dependency services
	# (database app) are running and healthy before the dependent service can come up.
	# docker run --rm -it -e DATABASE_URL=$(DATABASE_URL) -p $(HOST_PORT):$(CONTAINER_PORT) $(IMAGE_NAME):$(IMAGE_TAG)
	docker compose up -d backend

generate-mocks:
	mockgen -destination=mocks/mock_store.go -package=mocks github.com/swagnikdutta/one2n-sre-bootcamp/student Store

