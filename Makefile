
all: ensure test build

ensure:
	dep ensure

test: ensure
	go test -v -coverprofile=coverage.out $$(go list ./... | grep -v '/vendor/') && go tool cover -func=coverage.out

build: ensure
	go build -ldflags "-X gitlab.com/yakshaving.art/propaganda/version.Version=$(VERSION) -X gitlab.com/yakshaving.art/propaganda/version.Commit=$(COMMIT_ID) -X gitlab.com/yakshaving.art/propaganda/version.Date=$(COMMIT_DATE)"
