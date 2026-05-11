package main

import (
	"os"
	"syscall"
)

func newSocket() (*os.File, *os.File) {

	fd, err := syscall.Socketpair(
		syscall.AF_UNIX,
		syscall.SOCK_STREAM,
		0,
	)
	if err != nil {
		panic(err)
	}

	parentSock := os.NewFile(uintptr(fd[0]), "parentSock")
	childSock := os.NewFile(uintptr(fd[1]), "childSock")

	return parentSock, childSock
}