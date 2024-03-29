# go-cmdline

## Introduction
`cmdline` is a Go library to parse command line options (with optional default
values), arguments and subcommands.

## Usage
The following example is the simplest `cmdline` application possible:

```go
package main

import (
	"os"

	"github.com/galdor/go-cmdline"
)

func main() {
	cl := cmdline.New()
	cl.Parse(os.Args)
}
```

The resulting application handles `-h` and `--help`.

The `examples` directory contains examples for the various features of
`cmdline`. You can run them with `go run`. Feel free to copy and use these
examples in your own application.

## Contact
If you have an idea or a question, email me at <khaelin@gmail.com>.
