package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func nixEnv(path string) ([]Item, error) {
	cmd := Command("nix-env", "-f", path, "--drv-path", "-qaP", "*", "--xml", "--meta")
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

func missingPackages(path string) (map[string]bool, error) {
	var out bytes.Buffer
	cmd := Command("nix-build", "--dry-run", path)
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

func nixBuild(path string, attrs []string) error {
	buildArgs := []string{path}
	for _, attr := range attrs {
		buildArgs = append(buildArgs, "-k", "-A", attr)
	}

	cmd := Command("nix-build", buildArgs...)
	return cmd.Run()
}

func die(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		die("USAGE: %s path\n", args[0])
	}
	path := args[1]
	items, err := nixEnv(path)
	if err != nil {
		die("%s", err)
	}
	missingDrvs, err := missingPackages(path)
	if err != nil {
		die("%s", err)
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
	if err := nixBuild(path, missingAttrs); err != nil {
		die("nix-build failed: %s", err)
	}
}
