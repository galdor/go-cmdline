// Copyright (c) 2016-2019 Nicolas Martyanoff <khaelin@gmail.com>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package main

import (
	"fmt"
	"os"

	"github.com/galdor/go-cmdline"
)

func main() {
	cmdline := cmdline.New()

	cmdline.AddCommand("foo", "subcommand 1")
	cmdline.AddCommand("bar", "subcommand 2")

	cmdline.Parse(os.Args)

	var cmdFun func([]string)

	switch cmdline.CommandName() {
	case "foo":
		cmdFun = CmdFoo
	case "bar":
		cmdFun = CmdBar
	}

	cmdFun(cmdline.CommandNameAndArguments())
}

func CmdFoo(args []string) {
	fmt.Printf("running command \"foo\" with arguments %v\n", args[1:])
}

func CmdBar(args []string) {
	cmdline := cmdline.New()
	cmdline.AddOption("n", "", "value", "an example value")
	cmdline.Parse(args)

	fmt.Printf("running command \"bar\" with arguments %v\n", args[1:])

	if cmdline.IsOptionSet("n") {
		fmt.Printf("n: %s\n", cmdline.OptionValue("n"))
	}
}
