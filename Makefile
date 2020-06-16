REMOTE_DEPS = go.mod go.sum
GO_FILES = $(shell find . -name "*.go" | grep -v "pb.go$$" | grep -v /vendor/ | xargs)

function.zip: main
	zip function.zip main

main: $(GO_FILES) $(REMOTE_DEPS) types_easyjson.go cmd/bindata.go
	GOOS=linux go build ./cmd/main.go ./cmd/bindata.go

types_easyjson.go: types.go
	go get -u github.com/mailru/easyjson/...
	easyjson -all types.go

cmd/bindata.go: config.yaml
	go get -u github.com/go-bindata/go-bindata/...
	go-bindata -o cmd/bindata.go config.yaml
