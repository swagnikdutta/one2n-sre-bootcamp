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

generate-mocks:
