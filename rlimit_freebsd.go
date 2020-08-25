// +build !darwin

package main

import "syscall"

func rlimitMax(rlimit syscall.Rlimit) int64 {
	return rlimit.Max
}
