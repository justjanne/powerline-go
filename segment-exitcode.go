package main

import (
	"fmt"
	pwl "github.com/justjanne/powerline-go/powerline"
	"strconv"
	"syscall"
)

func getMeaningFromExitCode(exitCode int) string {
	switch exitCode {
	case 1:
		return "ERROR"
	case 2:
		return "USAGE"
	case 126:
		return "NOPERM"
	case 127:
		return "NOTFOUND"
	case 128 + int(syscall.SIGHUP):
		return "SIGHUP"
	case 128 + int(syscall.SIGINT):
		return "SIGINT"
	case 128 + int(syscall.SIGQUIT):
		return "SIGQUIT"
	case 128 + int(syscall.SIGILL):
		return "SIGILL"
	case 128 + int(syscall.SIGTRAP):
		return "SIGTRAP"
	case 128 + int(syscall.SIGABRT):
		return "SIGABRT"
	case 128 + int(syscall.SIGBUS):
		return "SIGBUS"
	case 128 + int(syscall.SIGFPE):
		return "SIGFPE"
	case 128 + int(syscall.SIGKILL):
		return "SIGKILL"
	case 128 + int(syscall.SIGUSR1):
		return "SIGUSR1"
	case 128 + int(syscall.SIGSEGV):
		return "SIGSEGV"
	case 128 + int(syscall.SIGUSR2):
		return "SIGUSR2"
	case 128 + int(syscall.SIGPIPE):
		return "SIGPIPE"
	case 128 + int(syscall.SIGALRM):
		return "SIGALRM"
	case 128 + int(syscall.SIGTERM):
		return "SIGTERM"
	case 128 + int(syscall.SIGSTKFLT):
		return "SIGSTKFLT"
	case 128 + int(syscall.SIGCHLD):
		return "SIGCHLD"
	case 128 + int(syscall.SIGCONT):
		return "SIGCONT"
	case 128 + int(syscall.SIGSTOP):
		return "SIGSTOP"
	case 128 + int(syscall.SIGTSTP):
		return "SIGTSTP"
	case 128 + int(syscall.SIGTTIN):
		return "SIGTTIN"
	case 128 + int(syscall.SIGTTOU):
		return "SIGTTOU"
	case 128 + int(syscall.SIGURG):
		return "SIGURG"
	case 128 + int(syscall.SIGXCPU):
		return "SIGXCPU"
	case 128 + int(syscall.SIGXFSZ):
		return "SIGXFSZ"
	case 128 + int(syscall.SIGVTALRM):
		return "SIGVTALRM"
	case 128 + int(syscall.SIGPROF):
		return "SIGPROF"
	case 128 + int(syscall.SIGWINCH):
		return "SIGWINCH"
	case 128 + int(syscall.SIGIO):
		return "SIGIO"
	case 128 + int(syscall.SIGPWR):
		return "SIGPWR"
	case 128 + int(syscall.SIGSYS):
		return "SIGSYS"
	default:
		return fmt.Sprintf("%d", exitCode)
	}
}

func segmentExitCode(p *powerline) []pwl.Segment {
	var meaning string
	if *p.args.PrevError != 0 {
		if *p.args.NumericExitCodes {
			meaning = strconv.Itoa(*p.args.PrevError)
		} else {
			meaning = getMeaningFromExitCode(*p.args.PrevError)
		}
		return []pwl.Segment{{
			Name:       "exit",
			Content:    meaning,
			Foreground: p.theme.CmdFailedFg,
			Background: p.theme.CmdFailedBg,
		}}
	}
	return []pwl.Segment{}
}
