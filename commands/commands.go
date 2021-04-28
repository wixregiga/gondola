package commands

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"gondola/app"
	"gondola/internal/runtimeutil"
	"gondola/util/stringutil"
)

var (
	commands = map[string]*command{}
	executed = false
)

type command struct {
	handler app.Handler
	opts    options
}

// Register registers a new command with the
// given function and options.
func Register(f app.Handler, opts ...Option) error {
	o := options{}
	for _, f := range opts {
		o = f(o)
	}
	if o.Name == "" {
		qname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		p := strings.Split(qname, ".")
		o.Name = p[len(p)-1]
		if o.Name == "" {
			return fmt.Errorf("could not determine name for function %v. Please, provide a name using Name().", f)
		}
		o.Name = stringutil.CamelCaseToLower(o.Name, "-")
	}
	if _, ok := commands[o.Name]; ok {
		return fmt.Errorf("duplicate command name %q", o.Name)
	}
	commands[o.Name] = &command{
		handler: f,
		opts:    o,
	}
	return nil
}

// Remove eliminates a previously registered command.
func Remove(name string) {
	for k := range commands {
		if k == name {
			delete(commands, k)
			break
		}
	}
}

// MustRegister works like Register, but panics
// if there's an error
func MustRegister(f app.Handler, opts ...Option) {
	if err := Register(f, opts...); err != nil {
		panic(err)
	}
}

// Error stops the command and prints the
// given error.
func Error(args ...interface{}) {
	panic(fmt.Sprint(args...))
}

// Errorf works like Error, but accepts a format parameter.
func Errorf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

// UsageError stops the command and prints the
// given error followed by the command usage.
func UsageError(args ...interface{}) {
	usageErrors(fmt.Sprint(args...))
}

// UsageErrorf works like UsageError, but accepts
// a format parameter.
func UsageErrorf(format string, args ...interface{}) {
	usageErrors(fmt.Sprintf(format, args...))
}

type usageError string

func (e usageError) Error() string {
	return string(e)
}

func usageErrors(s string) {
	err := usageError(s)
	panic(err)
}

func executeCommand(name string, cmd *command, args []string, a *app.App) (err error) {
	// Parse command flags
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	set.Usage = func() {
		commandHelp(name, -1, os.Stderr)
	}
	flags := map[string]interface{}{}
	for _, arg := range cmd.opts.Flags {
		switch arg.typ {
		case typBool:
			var b bool
			set.BoolVar(&b, arg.name, arg.def.(bool), arg.help)
			flags[arg.name] = &b
		case typInt:
			var i int
			set.IntVar(&i, arg.name, arg.def.(int), arg.help)
			flags[arg.name] = &i
		case typString:
			var s string
			set.StringVar(&s, arg.name, arg.def.(string), arg.help)
			flags[arg.name] = &s
		default:
			panic("invalid arg type")
		}
	}
	// Print error/help messages ourselves
	set.SetOutput(ioutil.Discard)
	err = set.Parse(args)
	if err != nil {
		if err == flag.ErrHelp {
			return
		}
		if strings.Contains(err.Error(), "provided but not defined") {
			flagName := strings.TrimSpace(strings.Split(err.Error(), ":")[1])
			fmt.Fprintf(os.Stderr, "command %s does not accept flag %s\n", name, flagName)
			return
		}
		return err
	}
	var params []string
	paramValues := make(map[string]string)
	for _, arg := range cmd.opts.Flags {
		params = append(params, arg.name)
		paramValues[arg.name] = fmt.Sprintf("%v", reflect.ValueOf(flags[arg.name]).Elem().Interface())
	}
	provider := &contextProvider{
		args:        set.Args(),
		params:      params,
		paramValues: paramValues,
	}
	if os.Getenv("GONDOLA_COMMANDS_NO_RECOVER") == "" {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("%v", r)
				}
				if file, line, ok := runtimeutil.PanicLocation(); ok {
					err = fmt.Errorf("%v (at %s:%d)", err, file, line)
				}
			}
		}()
	}
	ctx := a.NewContext(provider)
	defer a.CloseContext(ctx)
	cmd.handler(ctx)
	return
}

// Execute tries to run a command
// reading the parameters from the command line. It returs
// true if a command was executed and false if it wasn't.
// Note that most users won't need to call this function
// directly, since gndl.la/app.App will automatically call
// it before listening (and exit after executing the command
// if it was provided).
func Execute(a *app.App) (bool, error) {
	executed = true
	if !flag.Parsed() {
		flag.Parse()
	}
	args := flag.Args()
	if len(args) > 0 {
		cmd := strings.ToLower(args[0])
		for k, v := range commands {
			if cmd == k {
				if err := executeCommand(k, v, args[1:], a); err != nil {
					fmt.Fprintf(os.Stderr, "error running command %s: %s\n", cmd, err)
					if _, ok := err.(usageError); ok {
						commandHelp(cmd, -1, os.Stderr)
					}
				}
				return true, nil
			}
		}
		// Argument was given but it's not a command
		return false, fmt.Errorf("%s is not a registered command", args[0])
	}
	return false, nil
}

func execute(a *app.App) {
	if executed {
		return
	}
	done, err := Execute(a)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if done {
		os.Exit(0)
	}
}

// commandHelp prints the help for the given command
// to the given io.Writer
func commandHelp(name string, maxLen int, w io.Writer) {
	if maxLen < 0 {
		maxLen = len(name) + 1
	}
	fmt.Fprintf(w, "%s:%s%s\n", name, strings.Repeat(" ", maxLen-len(name)), commands[name].opts.Help)
	indent := strings.Repeat(" ", maxLen+1)
	if usage := commands[name].opts.Usage; usage != "" {
		fmt.Fprintf(w, "\n%sUsage: %s %s %s\n", indent, os.Args[0], name, usage)
	}
	if flags := commands[name].opts.Flags; len(flags) > 0 {
		fmt.Fprintf(w, "\n%sAvailable flags for %v:\n", indent, name)
		maxArgLen := -1
		helps := make([]string, len(flags))
		for ii, f := range flags {
			var buf bytes.Buffer
			buf.WriteByte('-')
			buf.WriteString(f.name)
			buf.WriteByte('=')
			if f.typ == typString {
				buf.WriteString(fmt.Sprintf("%q", f.def))
			} else {
				buf.WriteString(fmt.Sprintf("%v", f.def))
			}
			s := buf.String()
			if sl := len(s); sl > maxArgLen {
				maxArgLen = sl
			}
			helps[ii] = s
		}
		maxArgLen++
		format := fmt.Sprintf("%% -%ds", maxArgLen)
		for ii, f := range flags {
			fmt.Fprintf(w, indent)
			fmt.Fprintf(w, format, helps[ii])
			if f.help != "" {
				fmt.Fprintf(w, f.help)
			}
			fmt.Fprintf(w, "\n")
		}
	}
}

// commandsHelp prints the help for all commands to the given io.Writer
func commandsHelp(w io.Writer) {
	var cmds []string
	maxLen := 0
	for k, _ := range commands {
		if commandIsHidden(k) {
			continue
		}
		if l := len(k); l > maxLen {
			maxLen = l
		}
		cmds = append(cmds, k)
	}
	maxLen += 1
	sort.Strings(cmds)
	for _, v := range cmds {
		commandHelp(v, maxLen, w)
		fmt.Fprint(w, "\n\n")
	}
}

// Implementation of the help command for Gondola apps
func help(ctx *app.Context) {
	var cmd string
	ctx.ParseIndexValue(0, &cmd)
	if cmd != "" {
		c := strings.ToLower(cmd)
		if _, ok := commands[c]; ok {
			fmt.Fprintf(os.Stderr, "Help for command %s:\n", c)
			commandHelp(c, -1, os.Stderr)
		} else {
			fmt.Fprintf(os.Stderr, "No such command %q\n", cmd)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Commands:\n")
		commandsHelp(os.Stderr)
	}
}

func commandIsHidden(name string) bool {
	return name == "" || name[0] == '_'
}

func init() {
	MustRegister(help, Help("Show available commands with their respective help."))
	app.Signals.WillPrepare.Listen(execute)
}
