build:
	go build -o main queueListener.go

run:
	go run queueListener.go

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
