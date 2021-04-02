# healthcheck

A very simple Docker Health Check helper.

## Description

This small project can be used in Docker files to monitor the Docker container health by checking
an HTTP GET endpoint. If the endpoint is not responsive or returns an HTTP code < 200, or an HTTP code > 299,
the healthcheck executable (*hc*) will end with exit code 1. Otherwise, exit code 0 is used.

## Usage

Download the executable from [here](https://git.prolicht.digital/pub/healthcheck/-/releases/v1.0.0/downloads/binaries/hc).

Binary usage:
```bash
hc <url to check>

 - url to check: the full HTTP GET endpoint address that will be queried
 - exit code: 0 on success, 1 on failure
```


Docker usage:
```Dockerfile
COPY --from=buildImage /build/hc ./hc

HEALTHCHECK --interval=1s --timeout=1s --start-period=2s --retries=3 CMD [ "/hc", "www.url.to.check.com" ]
```

## Health check service library
This project also provides a very simple health check library that can be used to set up the health check endpoint.

This is an example of how to use the library:
```go
package main

import "git.prolicht.digital/pub/healthcheck"

func main() {
	healthcheck.New().Start() // start a default health check service on port 11223

	healthcheck.New(healthcheck.WithCustomCheck(func() int {
		return 2
	})).Start() // start a custom health check service

	healthcheck.New(healthcheck.ListenOn(":8080"), healthcheck.WithCustomCheck(func() int {
		return 2
	})).Start() // start a custom health check service on port 8080
}
```