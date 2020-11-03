Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X awssh/cmd.Version=$(Version) -X awssh/cmd.GitCommit=$(GitCommit)"
OUTDIR := bin

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

.PHONY: dist
dist:
	mkdir -p $(OUTDIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh-darwin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh-arm64
	CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh.exe