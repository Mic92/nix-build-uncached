package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// FIXME replace by sysconf?
const MAX_CHARS = 32 * 1024

func Command(cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	fmt.Printf("$ %s", cmd)
	for i, arg := range args {
		if i == 50 {
			fmt.Printf("...")
			break
		}
		fmt.Printf(" %s", arg)
	}
	fmt.Println()
	return c
}

func parseMissingDrvs(output *bytes.Buffer) map[string]bool {
	fmt.Println(output.String())
	scanner := bufio.NewScanner(output)
	scanner.Split(bufio.ScanLines)

	found := false
	missingDrvs := make(map[string]bool)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "these paths will be fetched") || strings.HasPrefix(line, "don't know how to build these paths") {
			break
		}
		if strings.HasPrefix(line, "these derivations will be built:") {
			found = true
		} else if found {
			drv := strings.TrimLeft(line, " ")
			missingDrvs[drv] = true
		}
	}

	return missingDrvs
}

func needExperimentalFlags() bool {
	cmd := exec.Command("nix")
	return cmd.Run() != nil
}

func nixDryBuild(buildArgs []string) (map[string]bool, error) {
	var out bytes.Buffer
	args := append([]string{"--dry-run"}, buildArgs...)
	cmd := Command("nix-build", args...)
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprint(os.Stderr, out.String())
		return nil, err
	}

	return parseMissingDrvs(&out), nil
}

func raiseFdLimit() (uint64, error) {
	var rlimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		return 0, fmt.Errorf("failed to get rlimit: %s", err)
	}

	if rlimit.Cur < rlimitMax(rlimit) {
		oldVal := rlimit.Cur
		rlimit.Cur = rlimitMax(rlimit)
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed setting rlimit: %s", err)
			return uint64(oldVal), nil
		}
	}
	return uint64(rlimit.Cur), nil
}

func nixBuild(drvs map[string]bool, buildArgs []string) error {
	buildArgs = append([]string{"build"}, buildArgs...)
	numBuildChars := len("nix") + 1
	for _, arg := range buildArgs {
		numBuildChars += len(arg) + 1
	}
	numChars := numBuildChars
	args := buildArgs

	if numChars > MAX_CHARS {
		return fmt.Errorf("too many arguments")
	}

	fdLimit, err := raiseFdLimit()
	if err != nil {
		return err
	}

	// nix build needs 3 fds per derivation, also add a safety margin on top.
	maxConcurrentBuilds := 1024 + fdLimit*3
	for drv := range drvs {
		n := len(drv) + 1
		if n+numChars > MAX_CHARS || uint64(len(args)-len(buildArgs)) >= maxConcurrentBuilds {
			cmd := Command("nix", args...)
			if err := cmd.Run(); err != nil {
				return err
			}
			numChars = numBuildChars
			args = buildArgs
		}
		args = append(args, drv)
		numChars += n
	}
	if numChars > numBuildChars {
		cmd := Command("nix", args...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func buildUncached(installables []string, buildArgs []string) ([]string, error) {
	if needExperimentalFlags() {
		buildArgs = append(buildArgs, "--experimental-features", "nix-command")
	}

	missingDrvs, err := nixDryBuild(append(installables, buildArgs...))
	if err != nil {
		return nil, fmt.Errorf("--dry-run failed: %s", err)
	}

	if err := nixBuild(missingDrvs, buildArgs); err != nil {
		return nil, fmt.Errorf("nix build failed: %s\n", err)
	}

	var builtDrvs []string
	for drv := range missingDrvs {
		builtDrvs = append(builtDrvs, drv)
	}

	return builtDrvs, nil
}
