# stage 1
FROM golang:1.24-alpine as builder
ARG BINARY_NAME=studentstore-bin
# when the image is built from the Makefile, the value of BINARY_NAME will be set to the value defined in the Makefile.
WORKDIR /app
COPY . .
# CGO_ENABLED=0 makes the binary runner a statically linked binary.
# When a program is statically linked, all the libraries it needs to run are present inside the binary itself.
# This is must if the base image is alpine, because alpine is lean. It does not have common system libraries like glibc.
# If the binary(runner) does not contain its library dependencies within itself, it will try to look for those dependencies
# in the alpine runtime — and it won't find it. Because alpine runtime won't have the C libraries.
#
# In our case, we are using debain as base image, so whether we have a statically linked binary or dynamically linked
# binary, it does not matter. But, static linking is still a nice safety check.
#RUN CGO_ENABLED=0 GOOS=linux go build -o appbin ./cmd

# I am using go-sqlite3, which is a cgo package. If you want to build your app using go-sqlite3 you need a gcc compiler
# present in your path, which is why CGO_ENABLED should be 1
RUN apk add --no-cache gcc musl-dev
# or use build-base for a complete set of tools — gcc, musl-dev, libc-dev, make, binutils and other core tools
# RUN apk add --no-cache build-base
RUN CGO_ENABLED=1 GOOS=linux go build -o ${BINARY_NAME} ./cmd


## stage 2
FROM alpine:3.22
WORKDIR /app
EXPOSE 8000
COPY --from=builder /app/${BINARY_NAME} .
CMD ["./studentstore-bin"]

