// Copyright (c) 2016,2017 Nicolas Martyanoff <khaelin@gmail.com>
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

package cmdline

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
)

// ---------------------------------------------------------------------------
//  Option
// ---------------------------------------------------------------------------
type Option struct {
	ShortName   string
	LongName    string
	ValueString string
	Description string
	Default     string

	Set   bool
	Value string
}

func (opt *Option) SortKey() string {
	if opt.ShortName != "" {
		return opt.ShortName
	}
	if opt.LongName != "" {
		return opt.LongName
	}
	return ""
}

type OptionArray []*Option

func (opts OptionArray) Len() int {
	return len(opts)
}

func (opts OptionArray) Less(i, j int) bool {
	return opts[i].SortKey() < opts[j].SortKey()
}

func (opts OptionArray) Swap(i, j int) {
	opts[i], opts[j] = opts[j], opts[i]
}

// ---------------------------------------------------------------------------
//  Argument
// ---------------------------------------------------------------------------
type Argument struct {
	Name        string
	Description string
	Trailing    bool

	Value          string
	TrailingValues []string
}

// ---------------------------------------------------------------------------
//  Command
// ---------------------------------------------------------------------------
type Command struct {
	Name        string
	Description string
}

// ---------------------------------------------------------------------------
//  Command line
// ---------------------------------------------------------------------------
type CmdLine struct {
	Options   map[string]*Option
	Arguments []*Argument

	Commands         map[string]*Command
	Command          string
	CommandArguments []string

	ProgramName string
}

func New() *CmdLine {
	c := &CmdLine{
		Options: make(map[string]*Option),

		Commands: make(map[string]*Command),
	}

	c.AddFlag("h", "help", "print help and exit")

	return c
}

func (c *CmdLine) AddFlag(short, long, desc string) {
	opt := &Option{
		ShortName:   short,
		LongName:    long,
		ValueString: "",
		Description: desc,
	}

	c.addOption(opt)
}

func (c *CmdLine) AddOption(short, long, value, desc string) {
	opt := &Option{
		ShortName:   short,
		LongName:    long,
		ValueString: value,
		Description: desc,
	}

	c.addOption(opt)
}

func (c *CmdLine) addOption(opt *Option) {
	if opt.ShortName != "" {
		if len(opt.ShortName) != 1 {
			panic("option short names must be one character long")
		}

		c.Options[opt.ShortName] = opt
	}

	if opt.LongName != "" {
		if len(opt.LongName) < 2 {
			panic("option long names must be at least two" +
				"characters long")
		}

		c.Options[opt.LongName] = opt
	}
}

func (c *CmdLine) SetOptionDefault(name, value string) {
	opt, found := c.Options[name]
	if !found {
		panic("unknown option")
	}

	if opt.ValueString == "" {
		panic("flags cannot have a default value")
	}

	opt.Default = value
}

func (c *CmdLine) AddArgument(name, desc string) {
	arg := &Argument{
		Name:        name,
		Description: desc,
	}

	c.addArgument(arg)
}

func (c *CmdLine) AddTrailingArguments(name, desc string) {
	arg := &Argument{
		Name:        name,
		Description: desc,
		Trailing:    true,
	}

	c.addArgument(arg)
}

func (c *CmdLine) addArgument(arg *Argument) {
	if len(c.Commands) > 0 {
		panic("cannot have both arguments and commands")
	}

	if len(c.Arguments) > 0 {
		last := c.Arguments[len(c.Arguments)-1]
		if last.Trailing {
			panic("cannot add argument after trailing argument")
		}
	}

	c.Arguments = append(c.Arguments, arg)
}

func (c *CmdLine) AddCommand(name, desc string) {
	if len(c.Arguments) == 0 {
		c.AddArgument("command", "the command to execute")
	} else if c.Arguments[0].Name != "command" {
		panic("cannot have both arguments and commands")
	}

	cmd := &Command{
		Name:        name,
		Description: desc,
	}

	c.Commands[cmd.Name] = cmd
}

func (c *CmdLine) Parse(args []string) {
	if len(args) == 0 {
		c.Die("empty argument array")
	}

	c.ProgramName = args[0]
	args = args[1:]

	for len(args) > 0 {
		arg := args[0]

		if arg == "--" {
			// End of options
			args = args[1:]
			break
		}

		isShort := len(arg) == 2 && arg[0] == '-' && arg[1] != '-'
		isLong := len(arg) > 2 && arg[0:2] == "--"

		if isShort || isLong {
			var key string
			if isShort {
				key = arg[1:2]
			} else {
				key = arg[2:]
			}

			opt, found := c.Options[key]
			if !found {
				c.Die("unknown option %q", key)
			}

			opt.Set = true

			if opt.ValueString == "" {
				args = args[1:]
			} else {
				if len(args) < 2 {
					c.Die("missing value "+
						"for option %q", key)
				}

				opt.Value = args[1]

				args = args[2:]
			}
		} else {
			// First argument
			break
		}
	}

	// Arguments
	if len(c.Arguments) > 0 && !c.IsOptionSet("help") {
		last := c.Arguments[len(c.Arguments)-1]

		min := len(c.Arguments)
		if last.Trailing {
			min--
		}

		if len(args) < min {
			c.Die("missing argument(s)")
		}

		for i := 0; i < min; i++ {
			c.Arguments[i].Value = args[i]
		}
		args = args[min:]

		if last.Trailing {
			last.TrailingValues = args
			args = args[len(args):]
		}
	}

	if len(c.Commands) > 0 {
		c.Command = c.Arguments[0].Value
		c.CommandArguments = args
	}

	if !c.IsOptionSet("help") {
		if len(c.Commands) > 0 {
			if _, found := c.Commands[c.Command]; !found {
				c.Die("unknown command %q", c.Command)
			}
		} else if len(args) > 0 {
			c.Die("invalid extra argument(s)")
		}
	}

	// Handle --help
	if c.IsOptionSet("help") {
		c.PrintUsage(os.Stdout)
		os.Exit(0)
	}
}

func (c *CmdLine) Die(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}

func (c *CmdLine) PrintUsage(w io.Writer) {
	// Header
	fmt.Fprintf(w, "Usage: %s OPTIONS", c.ProgramName)
	if len(c.Arguments) > 0 {
		for _, arg := range c.Arguments {
			if arg.Trailing {
				fmt.Fprintf(w, " [<%s> ...]", arg.Name)
			} else {
				fmt.Fprintf(w, " <%s>", arg.Name)
			}
		}
	}
	fmt.Fprintf(w, "\n\n")

	// Compute the width of the left column
	optStrs := make(map[*Option]string)
	maxWidth := 0

	for _, opt := range c.Options {
		if _, found := optStrs[opt]; found {
			continue
		}

		buf := bytes.NewBuffer([]byte{})

		if opt.ShortName != "" {
			fmt.Fprintf(buf, "-%s", opt.ShortName)
		}

		if opt.LongName != "" {
			if opt.ShortName != "" {
				buf.WriteString(", ")
			}

			fmt.Fprintf(buf, "--%s", opt.LongName)
		}

		if opt.ValueString != "" {
			fmt.Fprintf(buf, " <%s>", opt.ValueString)
		}

		str := buf.String()
		optStrs[opt] = str

		if len(str) > maxWidth {
			maxWidth = len(str)
		}
	}

	if len(c.Commands) > 0 {
		for name, _ := range c.Commands {
			if len(name) > maxWidth {
				maxWidth = len(name)
			}
		}
	} else if len(c.Arguments) > 0 {
		for _, arg := range c.Arguments {
			if len(arg.Name) > maxWidth {
				maxWidth = len(arg.Name)
			}
		}
	}

	// Print options
	fmt.Fprintf(w, "OPTIONS\n\n")

	opts := make([]*Option, len(optStrs))
	i := 0
	for opt, _ := range optStrs {
		opts[i] = opt
		i++
	}
	sort.Sort(OptionArray(opts))

	for _, opt := range opts {
		fmt.Fprintf(w, "%-*s  %s",
			maxWidth, optStrs[opt], opt.Description)

		if opt.Default != "" {
			fmt.Fprintf(w, " (default: %s)", opt.Default)
		}

		fmt.Fprintf(w, "\n")
	}

	if len(c.Commands) > 0 {
		// Print commands
		fmt.Fprintf(w, "\nCOMMANDS\n\n")
		for _, cmd := range c.Commands {
			fmt.Fprintf(w, "%-*s  %s\n",
				maxWidth, cmd.Name, cmd.Description)
		}
	} else if len(c.Arguments) > 0 {
		// Print arguments
		fmt.Fprintf(w, "\nARGUMENTS\n\n")

		for _, arg := range c.Arguments {
			fmt.Fprintf(w, "%-*s  %s\n",
				maxWidth, arg.Name, arg.Description)
		}
	}
}

func (c *CmdLine) IsOptionSet(name string) bool {
	opt, found := c.Options[name]
	if !found {
		panic("unknown option")
	}

	return opt.Set
}

func (c *CmdLine) OptionValue(name string) string {
	opt, found := c.Options[name]
	if !found {
		panic("unknown option")
	}

	if opt.Set {
		return opt.Value
	} else {
		return opt.Default
	}
}

func (c *CmdLine) ArgumentValue(name string) string {
	for _, arg := range c.Arguments {
		if arg.Name == name {
			return arg.Value
		}
	}

	panic("unknown argument")
}

func (c *CmdLine) TrailingArgumentsValues(name string) []string {
	if len(c.Arguments) == 0 {
		panic("empty argument array")
	}

	last := c.Arguments[len(c.Arguments)-1]
	if !last.Trailing {
		panic("no trailing arguments")
	}

	return last.TrailingValues
}

func (c *CmdLine) CommandName() string {
	if len(c.Commands) == 0 {
		panic("no command defined")
	}

	return c.Command
}

func (c *CmdLine) CommandArgumentsValues() []string {
	if len(c.Commands) == 0 {
		panic("no command defined")
	}

	return c.CommandArguments
}
