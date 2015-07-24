VERSION=$(shell git rev-parse --short HEAD)

NOW=$(shell date --iso-8601=seconds)

OUTPUT?='coduno'

all: test
	@go env
	@go version

	go build -x -v -o ${OUTPUT} \
		-ldflags "-X main.Version '${VERSION}' -X main.BuildTime '${NOW}'"

test: get
	go test -v

get: format
	go get -v

format:
	gofmt -d -e -l -s .
