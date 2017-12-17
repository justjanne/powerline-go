package main

import (
	"os"
)

func segmentShellVar(p *powerline) {
	var varContent string
	var varExists bool
	varContent, varExists = os.LookupEnv(*p.args.ShellVar);

	if (varExists) {
		p.appendSegment("shell-var", segment {
			content: "\\$" + *p.args.ShellVar + ":" + varContent,
			foreground: p.theme.ShellVarFg,
			background: p.theme.ShellVarBg,
		})
	}
}
