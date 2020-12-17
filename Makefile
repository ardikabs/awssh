Version := $(shell git describe --tags --dirty --always)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X awssh/cmd.Version=$(Version) -X awssh/cmd.GitCommit=$(GitCommit)"
OUTDIR := bin

GOLANGCI_VERSION = 1.31.0

.PHONY: test pretty mod tidy
test:
	CGO_ENABLED=0 go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover -coverprofile=cover.out

mod:
	go mod download

tidy:
	go mod tidy

pretty:
	gofmt -s -w **/*.go

build:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o awssh

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ./bin/ v${GOLANGCI_VERSION}
	@mv bin/golangci-lint "$@"

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run

.PHONY: dist
dist:
	mkdir -p $(OUTDIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh-darwin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh-arm64
	CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh.exe