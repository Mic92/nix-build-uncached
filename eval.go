package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
)

type Derivation struct {
	Outputs struct {
		Out struct {
			Path string `json:"path"`
		} `json:"out"`
	} `json:"outputs"`
}

type Item struct {
	AttrPath string `xml:"attrPath,attr"`
	DrvPath  string `xml:"drvPath,attr"`
	Outputs  []struct {
		Path string `xml:"path,attr"`
	} `xml:"output"`
}

type Items struct {
	XMLName xml.Name `xml:"items"`
	Items   []Item   `xml:"item"`
}

func nixEnv(path string, extraFlags []string) ([]Item, error) {
	args := []string{"-f", path, "--drv-path", "--out-path", "-qaP", "*", "--xml", "--meta"}
	args = append(args, extraFlags...)
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

func showDerivation(drvPath string) (map[string]Derivation, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("nix", "show-derivation", "--recursive", drvPath)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("nix show-derivation --recursive %s failed with: %s\n%s", drvPath, stderr.String(), err)
	}
	var derivations map[string]Derivation
	if err := json.Unmarshal(out.Bytes(), &derivations); err != nil {
		err := fmt.Errorf("Could not parse output of `nix show-derivation --recursive %s`: %s %s", drvPath, out.String(), stderr.String())
		return nil, err
	}
	return derivations, nil
}
