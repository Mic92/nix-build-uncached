package main

import (
	"fmt"
	"os"

	shellquote "github.com/kballard/go-shellquote"
)

type options struct {
	buildFlags   string
	installables []string
}

func parseFlags(args []string) (*options, error) {
	var opts options

	for i := 0; i < len(args); i++ {
		if args[i] == "-build-flags" {
			if i == len(args)-1 {
				return nil, fmt.Errorf("option '-build-flags' requires an argument")
			}
			i++
			opts.buildFlags = args[i]
			continue
		} else if args[i] == "-flags" {
			return nil, fmt.Errorf("option '-flags' is deprecated. You can now pass all those flags directly to nix-build-uncached")

		}
		opts.installables = append(opts.installables, args[i])
	}

	return &opts, nil
}

func realMain(args []string) error {
	opts, err := parseFlags(args)
	if err != nil {
		return err
	}

	buildFlags, err := shellquote.Split(opts.buildFlags)
	if err != nil {
		return fmt.Errorf("Value passed to -build-flags is not valid: %s", err)
	}

	if _, err := buildUncached(opts.installables, buildFlags); err != nil {
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
