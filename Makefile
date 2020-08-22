Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X awssh/internal/cli.Version=$(Version) -X awssh/internal/cli.GitCommit=$(GitCommit)"
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
	CGO_ENABLED=0 go build app/cli/main.go

.PHONY: dist
dist:
	mkdir -p $(OUTDIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh app/cli/main.go
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh-darwin app/cli/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh-armhf app/cli/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh-arm64 app/cli/main.go
	CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS) -o $(OUTDIR)/awssh.exe app/cli/main.go