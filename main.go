package main

import (
	"fmt"
	"os"

	shellquote "github.com/kballard/go-shellquote"
)

func realMain(args []string) error {
	var installables []string
	rawBuildFlags := "--keep-going"

	for i := 0; i < len(args); i++ {
		if args[i] == "-build-flags" {
			if i == len(args)-1 {
				return fmt.Errorf("option '-build-flags' requires an argument")
			}
			i++
			rawBuildFlags = args[i]
			continue
		}
		installables = append(installables, args[i])
	}

	buildFlags, err := shellquote.Split(rawBuildFlags)
	if err != nil {
		return fmt.Errorf("Value passed to -build-flags is not valid: %s", err)
	}

	if _, err := buildUncached(installables, buildFlags); err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}

func main() {
	if err := realMain(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}
