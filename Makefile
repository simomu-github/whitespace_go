
.PHONY: build deps clean

build: deps
	@GO111MODULE=on go build -o releases/ws cmd/ws.go

deps:
	@GO111MODULE=on go mod download

test: deps
	@GO111MODULE=on go test -v ./...

clean:
	rm releases/ws
