// Command slsa-l3-provenance-demo prints a greeting. It exists to show how
// release binaries can be built on GitHub Actions with SLSA v1.0 Build Level 3
// provenance; see README.md.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// version is set at release build time with -ldflags "-X main.version=...".
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("slsa-l3-provenance-demo", flag.ContinueOnError)
	flags.SetOutput(stderr)
	name := flags.String("name", "world", "who to greet")
	showVersion := flags.Bool("version", false, "print the version and exit")

	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() > 0 {
		fmt.Fprintf(stderr, "unexpected arguments: %q\n", flags.Args())
		flags.Usage()
		return 2
	}

	if *showVersion {
		fmt.Fprintln(stdout, version)
		return 0
	}
	fmt.Fprintln(stdout, Greeting(*name))
	return 0
}
