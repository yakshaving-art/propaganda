
# Variables
# UNAME		:= $(shell uname -s)

COMMIT_ID := `git log -1 --format=%H`
COMMIT_DATE := `git log -1 --format=%aI`
VERSION := $${CI_COMMIT_TAG:-SNAPSHOT-$(COMMIT_ID)}

.PHONY: ensure
ensure:
	@dep ensure

.PHONY: test
test: ensure
	go test -v -coverprofile=coverage.out $$(go list ./... | grep -v '/vendor/') && go tool cover -func=coverage.out

.PHONY: build
build: ensure
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X gitlab.com/yakshaving.art/propaganda/version.Version=$(VERSION) -X gitlab.com/yakshaving.art/propaganda/version.Commit=$(COMMIT_ID) -X gitlab.com/yakshaving.art/propaganda/version.Date=$(COMMIT_DATE)"
