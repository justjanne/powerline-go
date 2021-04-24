package main

import (
	"strings"
	"os/exec"
	pwl "github.com/justjanne/powerline-go/powerline"
)

func segmentFIB(p *powerline) []pwl.Segment {
	out, err := exec.Command("/sbin/sysctl", "-n", "net.my_fibnum").Output()

	if (err != nil) {
		return []pwl.Segment{}
	}

	var fib string = strings.TrimSpace(string(out))

	if (fib == "0") {
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "fib",
		Content:    p.symbols.FIBIndicator + " " + fib,
		Foreground: p.theme.HomeFg,
		Background: p.theme.HomeBg,
	}}
}
