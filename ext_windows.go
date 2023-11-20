//go:build windows

package main

import "syscall"

func reusePort(network, address string, conn syscall.RawConn) error {
	return conn.Control(func(descriptor uintptr) {
	})
}
