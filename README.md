# :dolphin: bkit ![GitHub Release](https://img.shields.io/github/v/release/nightnoryu/bkit) [![Build Status](https://github.com/nightnoryu/bkit/actions/workflows/check-go.yml/badge.svg)](https://github.com/nightnoryu/bkit/actions/workflows/check-go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/nightnoryu/bkit)](https://goreportcard.com/report/github.com/nightnoryu/bkit)

Container-native build system. Fork of [brewkit](https://github.com/ispringtech/brewkit).

The purpose of brewkit is full build process containerization, ensuring repeatable builds on any target machine (local or CI/CD). This fork is focused on adaptation of brewkit key features for open-source projects.

## Installation

Via `go install`:

```shell
go install github.com/nightnoryu/bkit/cmd/bkit@latest
```

Or just download a pre-build binary from the [releases page](https://github.com/nightnoryu/bkit/releases).

## Getting started

### Sample project

Let's create a simple bkit configuration for the following project structure:

```
bkit-example/
├── cmd
│   └── bkit-example
│       └── main.go
├── go.mod
├── go.sum
└── pkg
    └── functions
        ├── progressbar.go
        └── sayhi.go
```

```go
// sayhi.go

package functions

import "fmt"

func SayHi() {
        fmt.Println("Hi! Built with bkit")
}
```

```go
// progressbar.go

package functions

import (
	"time"

	"github.com/schollz/progressbar/v3"
)

func ShowProgressbar() {
	bar := progressbar.Default(100)
	for range 100 {
		bar.Add(1)
		time.Sleep(40 * time.Millisecond)
	}
}
```

```go
// main.go

package main

import "bkit-example/pkg/functions"

func main() {
	functions.SayHi()
	functions.ShowProgressbar()
}
```

In this example we're using Go modules and one external dependency.

### Sample configuration

To build this project using bkit, create a build configuration file `bkit.jsonnet`:

```jsonnet
local copy = std.native("copy");

local app = "bkit-example";

{
    apiVersion: "bkit/v1",
    targets: {
        all: ["build"],

        build: {
            from: "golang:1.24",
            workdir: "/app",
            copy: [
                copy("go.*", "."),
                copy("cmd", "cmd"),
                copy("pkg", "pkg"),
            ],
            command: "go build -o ./bin/" + app + " ./cmd/" + app,
            output: {
                artifact: "/app/bin/" + app,
                "local": "./bin"
            }
        }
    }
}
```

And you're set! Run the build by executing:

```shell
bkit build
```

### Breakdown

* In this example we created a single target to build a Go binary.
* bkit **targets** are similar to Makefile targets. You can also think of them as stages in a multi-stage Dockerfile.
* Targets are described in a format similar to the Dockerfile. For instance, in the `build` target we specify:
  * The docker image to run this target with
  * The working directory inside the container
  * Host files that we need to copy to the container
  * The command to run the build
  * And finally, the artifacts that we want to copy back to the host machine
* Executing `bkit build` runs the default target: `all`

For more examples and an in-depth documentation please refer to the [docs](docs).
