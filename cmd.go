package main

import (
	"fmt"
	"os"
	"os/exec"
)

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
