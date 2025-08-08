DB_URL=sqlite3://students.db
IMAGE_NAME=one2n
IMAGE_TAG=0.1.0
DB_PATH=students.db
HOST_PORT=8000
CONTAINER_PORT=8000

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
	migrate -database "$(DB_URL)" -path migrations up

migrate-down:
	migrate -database "$(DB_URL)" -path migrations down

migrate-force-drop:
	rm -f students.db

docker-build-alpine:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f alpine.Dockerfile .

docker-run:
	docker run --rm -it -e DB_PATH=$(DB_PATH) -p $(HOST_PORT):$(CONTAINER_PORT) $(IMAGE_NAME):$(IMAGE_TAG)

docker-clean:
	docker rmi $(IMAGE_NAME):$(IMAGE_TAG)

#generate-mocks:
#	mockgen -destination=mocks/mock_store.go -package=mocks github.com/swagnikdutta/one2n-sre-bootcamp/student Store

