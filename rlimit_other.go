// +build !darwin,!freebsd

package main

import "syscall"

func rlimitMax(rlimit syscall.Rlimit) uint64 {
	return rlimit.Max
}
