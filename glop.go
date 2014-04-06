package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/GeertJohan/go.linenoise"
	log "github.com/golang/glog"
	"github.com/palats/glop/parser"
	"github.com/palats/glop/runner"
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

	ctx := runner.NewContext(nil)
	for i := 0; ; i++ {
		s, err := linenoise.Line(fmt.Sprintf("In [%d]: ", i))

		if len(strings.TrimSpace(s)) > 0 {
			err := linenoise.AddHistory(s)
			if err != nil {
				log.Error(err)
			}
		}

		if err != nil {
			if err == linenoise.KillSignalError {
				return
			}
			fmt.Fprintf(os.Stderr, "Err [%d]: %s\n", i, err.Error())
		}

		result := parser.Parse(s).Eval(ctx)
		fmt.Printf("Out [%d]: <%T>%#+v\n", i, result, result)
	}
}
