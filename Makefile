all: build

build:
	go build -o bin/fb-messenger-analysis cmd/*.go

tools:
	go get -u github.com/golang/dep/cmd/dep

deps:
	dep ensure --vendor-only

clean:
	rm -rf bin