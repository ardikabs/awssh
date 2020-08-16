.PHONY: test pretty mod tidy

OUTDIR := deploy/_output

test:
	go test -coverprofile=cover.out ./...

mod:
	go mod download

tidy:
	go mod tidy

pretty:
	gofmt -s -w **/*.go

build-%:
	@GOOS=$* GOARCH=amd64 CGO_ENABLED=0 go build -o ${OUTDIR}/cli/awssh_amd64 app/cli/main.go