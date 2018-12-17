// Copyright (c) 2016-2018 Nicolas Martyanoff <khaelin@gmail.com>
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

	cmdline.AddArgument("foo", "the first argument")
	cmdline.AddArgument("bar", "the second argument")
	cmdline.AddTrailingArguments("name", "a trailing argument")

	cmdline.Parse(os.Args)

	fmt.Printf("foo: %s\n", cmdline.ArgumentValue("foo"))
	fmt.Printf("bar: %s\n", cmdline.ArgumentValue("bar"))
	fmt.Printf("names: %v\n", cmdline.TrailingArgumentsValues("name"))
}
