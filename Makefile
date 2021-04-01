PROGNAME=hc
GOCMD=go
GOFILES:=$(shell go list ./... | grep -v /vendor/)

# other targets

all: dep build

dep:
	GOPRIVATE=git.prolicht.digital/go/* $(GOCMD) mod download

build: dep
	GOPRIVATE=git.prolicht.digital/go/* GOOS=linux GOARCH=amd64 go build -o $(PROGNAME) -ldflags="-w -s" cmd/$(PROGNAME)/main.go
	sha256sum $(PROGNAME) >$(PROGNAME).sha256

clean: ## Remove build related file
	rm -fr ./$(PROGNAME)
	rm -fr ./$(PROGNAME).sha256