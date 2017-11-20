package interative

import (
	"os/user"
	"testing"
)

func TestExpandFilename(t *testing.T) {
	data := map[string]string{
		"/home/plo~p": "/home/plo~p",
		"~/plop":      "/foo/plop",
		// Unsupported expansion
		"~plop": "~plop",
	}

	for input, expected := range data {
		s, err := ExpandFilename(input, &user.User{HomeDir: "/foo"})
		if err != nil {
			t.Errorf("Unpexpected error: %v", err)
		}
		if s != expected {
			t.Errorf("Input %q - expected %q, got %q", input, expected, s)
		}
	}
}
