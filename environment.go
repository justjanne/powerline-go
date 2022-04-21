package main

import "runtime"

func homeEnvName() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	}
	return env
}
