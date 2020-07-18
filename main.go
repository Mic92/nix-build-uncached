package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"

	shellquote "github.com/kballard/go-shellquote"
)

func die(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func main() {
	flags := flag.String("flags", "", "additional arguments to pass to both nix-env/nix build")
	rawBuildFlags := flag.String("build-flags", "--keep-going", "additional arguments to pass to both nix build")

	flag.Parse()
	paths := flag.Args()
	if len(paths) != 1 {
		die("USAGE: %s path\n", os.Args[0])
	}
	path := paths[0]
	evalFlags, err := shellquote.Split(*flags)
	if err != nil {
		die("Value passed to -args is not valid: %s\n", err)
	}
	buildArgs, err := shellquote.Split(*rawBuildFlags)
	if err != nil {
		die("Value passed to -build-args is not valid: %s\n", err)
	}

	buildArgs = append(evalFlags, buildArgs...)

	builtItems, err := buildUncached(path, evalFlags, buildArgs)
	if err != nil {
		die("%s\n", err)
	}

	for _, item := range builtItems {
		if len(item.Outputs) == 0 {
			continue
		}
		output := item.Outputs[0]
		if !verifyPath(output.Path) {
			var out bytes.Buffer
			cmd := exec.Command("nix", "log", item.DrvPath)
			cmd.Stdout = &out
			if cmd.Run() == nil {
				fmt.Fprintf(os.Stderr, "%s could not be built:\n", item.AttrPath)
				fmt.Fprint(os.Stderr, out)
			} else {
				derivations, err := showDerivation(item.DrvPath)
			}
		}
	}
}
