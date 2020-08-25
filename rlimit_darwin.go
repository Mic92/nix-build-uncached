package main

import "syscall"

func rlimitMax(_ syscall.Rlimit) uint64 {
	// https://github.com/golang/go/issues/30401
	return 24576
}
