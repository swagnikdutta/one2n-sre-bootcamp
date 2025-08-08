# stage 1
FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 makes the binary runner a statically linked binary.
# When a program is statically linked, all the libraries it needs to run are present inside the binary itself.
# This is must if the base image is alpine, because alpine is lean. It does not have common system libraries like glibc.
# If the binary(runner) does not contain its library dependencies within itself, it will try to look for those dependencies
# in the alpine runtime — and it won't find it. Because alpine runtime won't have the C libraries — in that case the
# required libraries would have to be installed via apk (package manager of alpine)
#
# RUN CGO_ENABLED=0 GOOS=linux go build -o appbin ./cmd
#
# But since I am using go-sqlite3, which is a cgo package. If you want to build your app using go-sqlite3 you need a gcc compiler
# present in your path, which is why CGO_ENABLED should be 1 (written in documentation). Thus we install gcc and musl-dev (used
# for installing headers and static libraries) using apk (alpine's package manager). Or use build-base for a complete
# set of tools — gcc, musl-dev, libc-dev, make, binutils and other core tools.

RUN apk add --no-cache gcc musl-dev && \
    CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o main ./cmd
# Alternately use build-base. RUN apk add --no-cache build-base


## stage 2
FROM alpine:3.22
WORKDIR /app
EXPOSE 8000
COPY --from=builder /app/main .
CMD ["./main"]

