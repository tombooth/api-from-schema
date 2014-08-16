.PHONY: deps build clean

BINARY := api-from-schema
ORG_PATH := github.com/tombooth
REPO_PATH := $(ORG_PATH)/api-from-schema

all: deps build words

deps: third_party
	go run third_party.go get -t -v .

third_party:
	go run third_party.go setup $(REPO_PATH)

build:
	go run third_party.go build -v $(REPO_PATH)

words: words/post.html.template
	mkdir -p www
	pandoc --template=words/post.html.template -s words/article.md -o www/index.html
	cp -r words/img www/

words/post.html.template:
	curl https://raw.githubusercontent.com/tombooth/tombooth.github.io/develop/templates/post.html > words/post.html.template

release: clean all
	git checkout gh-pages && \
	git rm -rf . && \
	cp -r www/* . && \
	git add index.html && \
	git add -A img/ && \
	git commit -m "Update homepage" && \
	git push origin gh-pages -f && \
	git checkout master

clean:
	rm -rf third_party www
	rm -f api-from-schema
	rm -f words/post.html.template
