package main

var symbolTemplates = map[string]Symbols{
	"compatible": {
		Lock:          "RO",
		Network:       "SSH",
		Separator:     "\u25B6",
		SeparatorThin: "\u276F",
	},
	"patched": {
		Lock:          "\uE0A2",
		Network:       "\uE0A2",
		Separator:     "\uE0B0",
		SeparatorThin: "\uE0B1",

		RepoDetached:   "\u2693",
		RepoAhead:      "\u2B06",
		RepoBehind:     "\u2B07",
		RepoStaged:     "\u2714",
		RepoNotStaged:  "\u270E",
		RepoUntracked:  "+",
		RepoConflicted: "\u273C",
	},
	"flat": {
	},
}

var shellInfos = map[string]ShellInfo{
	"bash": {
		colorTemplate:    "\\[\\e%s\\]",
		rootIndicator:    " \\$ ",
		escapedBackslash: `\\\\`,
		escapedDollar:    `\$`,
	},
	"zsh": {
		colorTemplate:    "%%{\u001b%s%%}",
		rootIndicator:    " %# ",
		escapedBackslash: `\\`,
		escapedDollar:    `\$`,
	},
	"bare": {
		colorTemplate:    "%s",
		rootIndicator:    " $ ",
		escapedBackslash: `\`,
		escapedDollar:    `$`,
	},
}

var defaultTheme = Theme{
	Reset: 0xFF,

	UsernameFg:     250,
	UsernameBg:     240,
	UsernameRootBg: 124,

	HostnameFg: 250,
	HostnameBg: 238,

	HomeSpecialDisplay: true,
	HomeFg:             15,  // white
	HomeBg:             31,  // blueish
	PathFg:             250, // light grey
	PathBg:             237, // dark grey
	CwdFg:              254, // nearly-white grey
	SeparatorFg:        244,

	ReadonlyFg: 254,
	ReadonlyBg: 124,

	SshFg: 254,
	SshBg: 166, // medium orange

	DockerMachineFg: 177, // light purple
	DockerMachineBg: 55,  // purple

	RepoCleanFg: 0,   // black
	RepoCleanBg: 148, // a light green color
	RepoDirtyFg: 15,  // white
	RepoDirtyBg: 161, // pink/red

	JobsFg: 39,
	JobsBg: 238,

	CmdPassedFg: 15,
	CmdPassedBg: 236,
	CmdFailedFg: 15,
	CmdFailedBg: 161,

	SvnChangesFg: 22, // dark green
	SvnChangesBg: 148,

	GitAheadFg:      250,
	GitAheadBg:      240,
	GitBehindFg:     250,
	GitBehindBg:     240,
	GitStagedFg:     15,
	GitStagedBg:     22,
	GitNotStagedFg:  15,
	GitNotStagedBg:  130,
	GitUntrackedFg:  15,
	GitUntrackedBg:  52,
	GitConflictedFg: 15,
	GitConflictedBg: 9,

	VirtualEnvFg: 00,
	VirtualEnvBg: 35, // a mid-tone green
}
