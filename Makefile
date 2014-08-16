.PHONY: deps build clean

BINARY := api-from-schema
ORG_PATH := github.com/tombooth
REPO_PATH := $(ORG_PATH)/api-from-schema

all: deps build

deps: third_party
	go run third_party.go get -t -v .

third_party:
	go run third_party.go setup $(REPO_PATH)

build:
	go run third_party.go build -v $(REPO_PATH)

clean:
	rm -rf third_party
	rm api-from-schema
