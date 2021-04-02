PROGNAME=hc
GOCMD=go
GOFILES:=$(shell go list ./... | grep -v /vendor/)
MODULENAME=git.prolicht.digital/pub/healthcheck

all: dep test build

dep:
	$(GOCMD) mod download

build: dep
	GOOS=linux GOARCH=amd64 go build -o $(PROGNAME) -ldflags="-w -s" cmd/$(PROGNAME)/main.go
	sha256sum $(PROGNAME) >$(PROGNAME).sha256

test:
	$(GOCMD) test $(MODULENAME)/... -v -count=1

coverage:
	$(GOCMD) fmt $(GOFILES)
	$(GOCMD) test $(GOFILES) -v -coverprofile .testCoverage.txt
	$(GOCMD) tool cover -func=.testCoverage.txt  # use total:\s+\(statements\)\s+(\d+.\d+\%) as Gitlab CI regextotal:\s+\(statements\)\s+(\d+.\d+\%)

clean: ## Remove build related file
	rm -fr ./$(PROGNAME)
	rm -fr ./$(PROGNAME).sha256