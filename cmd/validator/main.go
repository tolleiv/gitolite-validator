package main

import (
	"io"
	"flag"
	"os"
	"fmt"
	validator "github.com/tolleiv/gitolite-validator"
)

func input(n int) io.ReadCloser {
	if flag.Arg(n) == "" {
		return os.Stdin
	} else {
		f, err := os.Open(flag.Arg(n))
		if err != nil {
			bail(1, "input error: %s", err)
		}
		return f
	}
}

func bail(status int, t string, args ...interface{}) {
	var w io.Writer
	if status == 0 {
		w = os.Stdout
	} else {
		w = os.Stderr
	}
	fmt.Fprintf(w, t + "\n", args...)
	os.Exit(status)
}

func main() {
	r := input(1)
	defer r.Close()

	if err := validator.Read(r); err != nil {
		bail(1, "parse error: %s", err)
	}
}
