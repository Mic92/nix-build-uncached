package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	shellquote "github.com/kballard/go-shellquote"
)

type Item struct {
	AttrPath string `xml:"attrPath,attr"`
	DrvPath  string `xml:"drvPath,attr"`
}

type Items struct {
	XMLName xml.Name `xml:"items"`
	Items   []Item   `xml:"item"`
}

func Command(cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	fmt.Printf("$ %s", cmd)
	for _, arg := range args {
		fmt.Printf(" %s", arg)
	}
	fmt.Println()
	return c
}

func nixEnv(path string, extraArgs []string) ([]Item, error) {
	args := []string{"-f", path, "--drv-path", "-qaP", "*", "--xml", "--meta"}
	args = append(args, extraArgs...)
	cmd := Command("nix-env", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", out.String())
		return nil, fmt.Errorf("nix-env failed: %v\n", err)
	}
	var items Items
	if err := xml.Unmarshal(out.Bytes(), &items); err != nil {
		return nil, fmt.Errorf("failed to parse nix-env output: %v", err)
	}
	return items.Items, nil
}

func missingPackages(path string, extraArgs []string) (map[string]bool, error) {
	var out bytes.Buffer
	args := []string{"--dry-run", path}
	args = append(args, extraArgs...)
	cmd := Command("nix-build", args...)
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
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
			fmt.Println(drv)
			missingDrvs[drv] = true
		}
	}

	return missingDrvs, nil

}

func nixBuild(path string, attrs []string, extraArgs []string) error {
	args := []string{"build", "-f", path}
	args = append(args, attrs...)
	args = append(args, extraArgs...)

	cmd := Command("nix", args...)
	return cmd.Run()
}

func die(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func buildUncached(path string, evalArgs []string, buildArgs []string) {
	items, err := nixEnv(path, evalArgs)
	if err != nil {
		die("%s\n", err)
	}

	missingDrvs, err := missingPackages(path, evalArgs)
	if err != nil {
		die("%s\n", err)
	}
	var missingAttrs []string

	for _, item := range items {
		if _, ok := missingDrvs[item.DrvPath]; ok {
			missingAttrs = append(missingAttrs, item.AttrPath)
		}
	}
	fmt.Printf("%d/%d attribute(s) will be built:\n", len(missingAttrs), len(items))
	for _, attr := range missingAttrs {
		fmt.Printf("  %s\n", attr)
	}
	if len(missingAttrs) == 0 {
		return
	}
	if err := nixBuild(path, missingAttrs, buildArgs); err != nil {
		die("nix-build failed: %s\n", err)
	}

}

func main() {
	args := flag.String("args", "", "additional arguments to pass to both nix-env/nix build")
	rawBuildArgs := flag.String("build-args", "--keep-going", "additional arguments to pass to both nix build")

	flag.Parse()
	paths := flag.Args()
	fmt.Println(paths)
	if len(paths) != 1 {
		die("USAGE: %s path\n", os.Args[0])
	}
	path := paths[0]
	evalArgs, err := shellquote.Split(*args)
	if err != nil {
		die("Value passed to -args is not valid: %s\n", err)
	}
	buildArgs, err := shellquote.Split(*rawBuildArgs)
	if err != nil {
		die("Value passed to -build-args is not valid: %s\n", err)
	}

	buildArgs = append(evalArgs, buildArgs...)

	buildUncached(path, evalArgs, buildArgs)
}
