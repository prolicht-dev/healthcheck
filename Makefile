PROGNAME=hc
GOCMD=go
GOFILES:=$(shell go list ./... | grep -v /vendor/)
MODULENAME=git.prolicht.digital/pub/healthcheck

all: dep test build

dep:
	$(GOCMD) mod download

build: dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(PROGNAME) -ldflags="-w -s -extldflags \"-static\"" cmd/$(PROGNAME)/main.go
	sha256sum $(PROGNAME) >$(PROGNAME).sha256

test:
	$(GOCMD) test $(MODULENAME)/... -v -count=1

coverage:
	$(GOCMD) fmt $(GOFILES)
	$(GOCMD) test $(GOFILES) -v -coverprofile .testCoverage.txt
	$(GOCMD) tool cover -func=.testCoverage.txt  # use total:\s+\(statements\)\s+(\d+.\d+\%) as Gitlab CI regextotal:\s+\(statements\)\s+(\d+.\d+\%)

coverage-html: coverage
	$(GOCMD) tool cover -html=.testCoverage.txt

clean: ## Remove build related file
	rm -fr ./$(PROGNAME)
	rm -fr ./$(PROGNAME).sha256
	rm -fr .testCoverage.txt