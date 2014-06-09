package repl

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/GeertJohan/go.linenoise"
	log "github.com/golang/glog"
	"github.com/palats/glop/parser"
	glopRuntime "github.com/palats/glop/runtime"
)

// ExpandFilename replace '~/' with the user home directory.
func ExpandFilename(s string, u *user.User) (string, error) {
	if strings.HasPrefix(s, "~/") {
		if u == nil {
			var err error
			u, err = user.Current()
			if err != nil {
				return "", err
			}
		}
		s = u.HomeDir + s[1:]
	}
	return s, nil
}

// REPL starts a full interpreter, accepting glop code on its prompt.
// If historyFilename is empty, no history is loaded/saved.
func REPL(historyFilename string) error {
	// Initialize history file
	if historyFilename != "" {
		fname, err := ExpandFilename(historyFilename, nil)
		if err != nil {
			return fmt.Errorf("%q is invalid: %v", historyFilename, err)
		}

		linenoise.LoadHistory(fname)
		defer linenoise.SaveHistory(fname)
	}

	// Prepare runtime environment
	ctx := glopRuntime.NewContext(nil)
	ctx.Set("exit", parser.NewAny(func() {
		os.Exit(0)
	}, nil))

	// REPL main loop
	for i := 0; ; i++ {
		// Prompt
		raw, err := linenoise.Line(fmt.Sprintf("In [%d]: ", i))
		if err == linenoise.KillSignalError {
			return nil
		}
		if len(strings.TrimSpace(raw)) == 0 {
			continue
		}

		if err := linenoise.AddHistory(raw); err != nil {
			log.Error(err)
		}

		// Parsing
		var out parser.Value
		root, err := parser.Parse(parser.NewSource(raw))
		if err == nil {
			// Executing
			out, err = ctx.SafeEval(root)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Error [%d]: %v\n", i, err))
			continue
		}

		v := out.Value()
		fmt.Printf("Out [%d]: <%T>%#+v\n", i, v, v)
	}
}
