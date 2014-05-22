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

func main() {
	flag.Parse()

	if *historyFilename != "" {
		fname, err := expandFilename(*historyFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%q is invalid: %v", *historyFilename, err)
			os.Exit(1)
		}

		linenoise.LoadHistory(fname)
		defer linenoise.SaveHistory(fname)
	}

	ctx := glopRuntime.NewContext(nil)
	ctx.Set("exit", func() {
		os.Exit(0)
	})

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

		tree, errs := parser.Parse(parser.NewSource(raw))
		if len(errs) > 0 {
			for _, err := range errs {
				prefix := fmt.Sprintf("Parse error [%d]: ", i)
				spaces := strings.Repeat(" ", len(prefix))
				fmt.Fprintf(os.Stderr, "%s %s\n", prefix, err.Error())

				if details, ok := err.(parser.Error); ok {
					line, err := details.Ref.Source.Line(details.Ref.Line)
					if err == nil {
						// TODO: This is probably going to end up corrupting the term if
						// the input is not clean, so we might want more escaping.
						fmt.Fprintf(os.Stderr, "%s %s\n", spaces, strings.TrimRight(string(line), "\n"))
						if details.Ref.Column >= 0 && details.Ref.Column <= len(line) {
							fmt.Fprintf(os.Stderr, "%s %s^\n", spaces, strings.Repeat("-", details.Ref.Column))
						}
					} else {
						fmt.Fprintf(os.Stderr, "%s unable to get source info: %s", spaces, err)
					}
				}
			}
			continue
		}

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
			fmt.Printf("Panic [%d]: %v\n", i, recov)
			for _, line := range strings.Split(string(stack), "\n") {
				fmt.Printf("  %s\n", line)
			}
		} else {
			fmt.Printf("Out [%d]: <%T>%#+v\n", i, result, result)
		}
	}
}
