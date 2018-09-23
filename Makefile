
all: ensure test build

ensure:
	dep ensure

test: ensure
	go test -v -coverprofile=coverage.out $$(go list ./... | grep -v '/vendor/') && go tool cover -func=coverage.out

build: ensure
	go build

snapshot:
	goreleaser --rm-dist --snapshot

release: ensure
	goreleaser --rm-dist
