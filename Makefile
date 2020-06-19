GOBIN = $(shell go env GOPATH)/bin
REMOTE_DEPS = go.mod go.sum
GO_FILES = $(shell find . -name "*.go" | grep -v "pb.go$$" | grep -v /vendor/ | xargs)

function.zip: main
	zip function.zip main

main: $(GO_FILES) $(REMOTE_DEPS) types_easyjson.go cmd/lambda/bindata.go
	GOOS=linux go build -ldflags="-s -w" ./cmd/lambda/main.go ./cmd/lambda/bindata.go

types_easyjson.go: types.go
	go get -u github.com/mailru/easyjson/...
	$(GOBIN)/easyjson -all types.go

cmd/lambda/bindata.go: config.yaml
	go get -u github.com/go-bindata/go-bindata/...
	$(GOBIN)/go-bindata -o cmd/lambda/bindata.go config.yaml
