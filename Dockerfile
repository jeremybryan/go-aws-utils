FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    QUEUE="" \
    PROFILE="" \
    REGION="" \
    THRESHOLD="" \
    GETEND=""

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main queueListener.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
#FROM scratch
FROM alpine:latest

COPY --from=builder /dist/main /

# Command to run
ENTRYPOINT /main -queue ${QUEUE:-""} -profile ${PROFILE:-"default"} -threshold ${THRESHOLD:-10} -region ${REGION:-""} -getEndpoint ${GETEND:-""}
