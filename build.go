package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func verifyPath(path string) bool {
	cmd := exec.Command("nix-store", "--verify-path", path)
	return cmd.Run() != nil
}

func missingPackages(path string, extraFlags []string) (map[string]bool, error) {
	var out bytes.Buffer
	args := []string{"--dry-run", path}
	args = append(args, extraFlags...)
	cmd := Command("nix-build", args...)
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprint(os.Stderr, out.String())
		return nil, err
	}

	scanner := bufio.NewScanner(&out)
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

	return missingDrvs, nil
}

func escapeAttr(attr string) string {
	parts := strings.Split(attr, ".")
	quoted := make([]string, len(parts))
	for idx, part := range parts {
		quoted[idx] = fmt.Sprintf("\"%s\"", part)
	}
	return strings.Join(quoted, ".")
}

func needExperimentalFlags() bool {
	cmd := exec.Command("nix")
	return cmd.Run() != nil
}

func nixBuild(path string, items []Item, extraArgs []string) error {
	args := []string{"build"}

	tmpFile, err := ioutil.TempFile("", "*.nix")
	if err != nil {
		die("failed to create temporary file: %s", err)
	}
	defer os.Remove(tmpFile.Name())
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		die("invalid '%s' passed", path)
	}
	header := fmt.Sprintf(`{...} @args:
let
  set' = import "%s";
  set = if builtins.isFunction set' then set' args else set';
in [
`, absolutePath)
	tmpFile.WriteString(header)
	for _, item := range items {
		tmpFile.WriteString(fmt.Sprintf("set.%s\n", escapeAttr(item.AttrPath)))
	}
	tmpFile.WriteString("]")
	tmpFile.Close()
	if needExperimentalFlags() {
		args = append(args, []string{"--experimental-features", "nix-command"}...)
	}
	args = append(args, []string{"-f", tmpFile.Name()}...)
	args = append(args, extraArgs...)
	cmd := Command("nix", args...)
	return cmd.Run()
}

func buildUncached(path string, evalFlags []string, buildArgs []string) ([]Item, error) {
	items, err := nixEnv(path, evalFlags)
	if err != nil {
		return nil, err
	}

	missingDrvs, err := missingPackages(path, evalFlags)
	if err != nil {
		return nil, err
	}
	var missingItems []Item

	for _, item := range items {
		if _, ok := missingDrvs[item.DrvPath]; ok {
			missingItems = append(missingItems, item)
			missingDrvs[item.DrvPath] = false
		}
	}

	fmt.Printf("%d/%d attribute(s) will be built:\n", len(missingItems), len(items))
	for _, attr := range missingItems {
		fmt.Printf("  %s\n", attr)
	}
	if len(missingItems) == 0 {
		return missingItems, nil
	}

	if err := nixBuild(path, missingItems, buildArgs); err != nil {
		return nil, fmt.Errorf("nix build failed: %s\n", err)
	}

	return missingItems, nil
}
