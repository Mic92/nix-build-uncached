// +build !darwin

package main

import "syscall"

func rlimitMax(rlimit syscall.Rlimit) uint64 {
	return uint64(rlimit.Max)
}
