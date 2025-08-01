APP_NAME=student-crud
DB_URL=sqlite3://students.db

run:
# just cmd does not work as go interprets it to be package path, containing multiple source files/sub-packages even
# if it doesn't. So we add ./ so that it is treated as an import path.
	go run ./cmd

build:
	go build -o $(APP_NAME) ./cmd

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f $(APP_NAME)

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

#generate-mocks:
#	mockgen -destination=mocks/mock_store.go -package=mocks github.com/swagnikdutta/one2n-sre-bootcamp/student Store

