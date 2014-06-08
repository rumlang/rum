package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/palats/glop/repl"
)

var historyFilename = flag.String(
	"history_filename", "~/.glophistory",
	"File to use for prompt history. '~/' is expanded to current user home "+
		"directory. Set to empty string to disable history loading/saving.")

func main() {
	flag.Parse()
	if err := repl.REPL(*historyFilename); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
