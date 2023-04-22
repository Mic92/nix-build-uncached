package main

import (
	"bytes"
	"fmt"
	"os"

	shellquote "github.com/kballard/go-shellquote"
)

type options struct {
	buildFlags   string
	installables []string
}

type nixVersion struct {
	Major, Minor, Patch uint64
}

func getNixVersion() (nixVersion, error) {
	var version nixVersion
	cmd := Command("nix", "--version")
	var outb bytes.Buffer
	cmd.Stdout = &outb
	if err := cmd.Run(); err != nil {
		return version, err
	}
	_, err := fmt.Sscanf(outb.String(), "nix (Nix) %d.%d.%d", &version.Major, &version.Minor, &version.Patch)
	if err != nil {
		return version, err
	}
	return version, nil
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
	version, err := getNixVersion()
	if err != nil {
		return fmt.Errorf("Failed to get nix version: %s", err)
	}

	if _, err := buildUncached(opts.installables, buildFlags, version); err != nil {
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
