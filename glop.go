package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/GeertJohan/go.linenoise"
	log "github.com/golang/glog"
	"github.com/palats/glop/parser"
	glopRuntime "github.com/palats/glop/runtime"
)

var historyFilename = flag.String(
	"history_filename", "~/.glophistory",
	"File to use for prompt history. '~/' is expanded to current user home "+
		"directory. Set to empty string to disable history loading/saving.")

// expandFilename replace '~/' with the user home directory.
func expandFilename(s string) (string, error) {
	if strings.HasPrefix(s, "~/") {
		u, err := user.Current()
		if err != nil {
			return "", err
		}
		s = u.HomeDir + s[1:]
	}
	return s, nil
}

// refContext prints on stderr what is on the line referenced by the provided
// source reference.
func refContext(ref parser.SourceRef, prefix string) {
	line, err := ref.Source.Line(ref.Line)
	if err == nil {
		// TODO: This is probably going to end up corrupting the term if
		// the input is not clean, so we might want more escaping.
		fmt.Fprintf(os.Stderr, "%s%s\n", prefix, strings.TrimRight(string(line), "\n"))
		if ref.Column >= 0 && ref.Column <= len(line) {
			fmt.Fprintf(os.Stderr, "%s%s^\n", prefix, strings.Repeat("-", ref.Column))
		}
	} else {
		fmt.Fprintf(os.Stderr, "%sunable to get source info: %s", prefix, err)
	}
}

func main() {
	flag.Parse()

	// Initialize history file
	if *historyFilename != "" {
		fname, err := expandFilename(*historyFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%q is invalid: %v", *historyFilename, err)
			os.Exit(1)
		}

		linenoise.LoadHistory(fname)
		defer linenoise.SaveHistory(fname)
	}

	// Prepare runtime environment
	ctx := glopRuntime.NewContext(nil)
	ctx.Set("exit", func() {
		os.Exit(0)
	})

	// REPL main loop
	for i := 0; ; i++ {
		raw, err := linenoise.Line(fmt.Sprintf("In [%d]: ", i))
		if err == linenoise.KillSignalError {
			return
		}
		if len(strings.TrimSpace(raw)) == 0 {
			continue
		}

		if err := linenoise.AddHistory(raw); err != nil {
			log.Error(err)
		}

		// Parsing
		tree, errs := parser.Parse(parser.NewSource(raw))
		if len(errs) > 0 {
			for _, err := range errs {
				prefix := fmt.Sprintf("Parse error [%d]: ", i)
				fmt.Fprintf(os.Stderr, "%s %s\n", prefix, err.Error())

				if details, ok := err.(parser.Error); ok {
					refContext(details.Ref, strings.Repeat(" ", len(prefix)+1))
				}
			}
			continue
		}

		// Execution
		var result, recov interface{}
		var stack []byte
		func() {
			defer func() {
				const size = 16384
				stack = make([]byte, size)
				// Unfortunately, that also catch itself, adding noise to the trace.
				stack = stack[:runtime.Stack(stack, false)]
				recov = recover()
			}()
			result = ctx.Eval(tree)
		}()
		if recov != nil {
			prefix := fmt.Sprintf("Panic [%d]: ", i)
			fmt.Fprintf(os.Stderr, "%s %v\n", prefix, recov)
			if details, ok := recov.(glopRuntime.Error); ok {
				refContext(details.Ref, strings.Repeat(" ", len(prefix)+1))
			} else {
				for _, line := range strings.Split(string(stack), "\n") {
					fmt.Printf("  %s\n", line)
				}
			}
		} else {
			fmt.Printf("Out [%d]: <%T>%#+v\n", i, result, result)
		}
	}
}
