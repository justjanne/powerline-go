// +build !windows

package main

import (
	"os"
)

func userIsAdmin() bool {
	return os.Getuid() == 0
}
