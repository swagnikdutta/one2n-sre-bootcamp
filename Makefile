APP_NAME := student-crud

run:
	go run .

build:
	go build -o $(APP_NAME) .

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

#generate-mocks:
#	mockgen -destination=mocks/mock_store.go -package=mocks github.com/swagnikdutta/one2n-sre-bootcamp/student Store

