package main

import (
	"flag"
	"strings"
)

type arguments struct {
	CwdMode                *string
	CwdMaxDepth            *int
	CwdMaxDirSize          *int
	ColorizeHostname       *bool
	HostnameOnlyIfSSH      *bool
	SshAlternateIcon       *bool
	EastAsianWidth         *bool
	PromptOnNewLine        *bool
	StaticPromptIndicator  *bool
	VenvNameSizeLimit      *int
	GitAssumeUnchangedSize *int64
	GitDisableStats        *string
	GitMode                *string
	Jobs                   *int
	Mode                   *string
	Theme                  *string
	Shell                  *string
	Modules                *string
	ModulesRight           *string
	Priority               *string
	MaxWidthPercentage     *int
	TruncateSegmentWidth   *int
	PrevError              *int
	NumericExitCodes       *bool
	IgnoreRepos            *string
	ShortenGKENames        *bool
	ShortenEKSNames        *bool
	ShortenOpenshiftNames  *bool
	ShellVar               *string
	ShellVarNoWarnEmpty    *bool
	TrimADDomain           *bool
	PathAliases            *string
	Duration               *string
	DurationMin            *string
	DurationLowPrecision   *bool
	Eval                   *bool
	Condensed              *bool
	IgnoreWarnings         *bool
	Time                   *string
	ViMode                 *string
}

var args = arguments{
	CwdMode: flag.String(
		"cwd-mode",
		defaults.CwdMode,
		commentsWithDefaults("How to display the current directory",
			"(valid choices: fancy, semifancy, plain, dironly)")),
	CwdMaxDepth: flag.Int(
		"cwd-max-depth",
		defaults.CwdMaxDepth,
		commentsWithDefaults("Maximum number of directories to show in path")),
	CwdMaxDirSize: flag.Int(
		"cwd-max-dir-size",
		defaults.CwdMaxDirSize,
		commentsWithDefaults("Maximum number of letters displayed for each directory in the path")),
	ColorizeHostname: flag.Bool(
		"colorize-hostname",
		defaults.ColorizeHostname,
		comments("Colorize the hostname based on a hash of itself, or use the PLGO_HOSTNAMEFG and PLGO_HOSTNAMEBG env vars (both need to be set).")),
	HostnameOnlyIfSSH: flag.Bool(
		"hostname-only-if-ssh",
		defaults.HostnameOnlyIfSSH,
		comments("Show hostname only for SSH connections")),
	SshAlternateIcon: flag.Bool(
		"alternate-ssh-icon",
		defaults.SshAlternateIcon,
		comments("Show the older, original icon for SSH connections")),
	EastAsianWidth: flag.Bool(
		"east-asian-width",
		defaults.EastAsianWidth,
		comments("Use East Asian Ambiguous Widths")),
	PromptOnNewLine: flag.Bool(
		"newline",
		defaults.PromptOnNewLine,
		comments("Show the prompt on a new line")),
	StaticPromptIndicator: flag.Bool(
		"static-prompt-indicator",
		defaults.StaticPromptIndicator,
		comments("Always show the prompt indicator with the default color, never with the error color")),
	VenvNameSizeLimit: flag.Int(
		"venv-name-size-limit",
		defaults.VenvNameSizeLimit,
		comments("Show indicator instead of virtualenv name if name is longer than this limit (defaults to 0, which is unlimited)")),
	Jobs: flag.Int(
		"jobs",
		defaults.Jobs,
		comments("Number of jobs currently running")),
	GitAssumeUnchangedSize: flag.Int64(
		"git-assume-unchanged-size",
		defaults.GitAssumeUnchangedSize,
		comments("Disable checking for changed/edited files in git repositories where the index is larger than this size (in KB), improves performance")),
	GitDisableStats: flag.String(
		"git-disable-stats",
		strings.Join(defaults.GitDisableStats, ","),
		commentsWithDefaults("Comma-separated list to disable individual git statuses",
			"(valid choices: ahead, behind, staged, notStaged, untracked, conflicted, stashed)")),
	GitMode: flag.String(
		"git-mode",
		defaults.GitMode,
		commentsWithDefaults("How to display git status",
			"(valid choices: fancy, compact, simple)")),
	Mode: flag.String(
		"mode",
		defaults.Mode,
		commentsWithDefaults("The characters used to make separators between segments.",
			"(valid choices: patched, compatible, flat)")),
	Theme: flag.String(
		"theme",
		defaults.Theme,
		commentsWithDefaults("Set this to the theme you want to use",
			"(valid choices: default, low-contrast, gruvbox, solarized-dark16, solarized-light16)")),
	Shell: flag.String(
		"shell",
		defaults.Shell,
		commentsWithDefaults("Set this to your shell type",
			"(valid choices: autodetect, bare, bash, zsh)")),
	Modules: flag.String(
		"modules",
		strings.Join(defaults.Modules, ","),
		commentsWithDefaults("The list of modules to load, separated by ','",
			"(valid choices: aws, bzr, cwd, direnv, docker, docker-context, dotenv, duration, exit, fossil, gcp, git, gitlite, goenv, hg, host, jobs, kube, load, newline, nix-shell, node, perlbrew, perms, plenv, rbenv, root, rvm, shell-var, shenv, ssh, svn, termtitle, terraform-workspace, time, user, venv, vgo, vi-mode, wsl)",
			"Unrecognized modules will be invoked as 'powerline-go-MODULE' executable plugins and should output a (possibly empty) list of JSON objects that unmarshal to powerline-go's Segment structs.")),
	ModulesRight: flag.String(
		"modules-right",
		strings.Join(defaults.ModulesRight, ","),
		comments("The list of modules to load anchored to the right, for shells that support it, separated by ','",
			"(valid choices: aws, bzr, cwd, direnv, docker, docker-context, dotenv, duration, exit, fossil, gcp, git, gitlite, goenv, hg, host, jobs, kube, load, newline, nix-shell, node, perlbrew, perms, plenv, rbenv, root, rvm, shell-var, shenv, ssh, svn, termtitle, terraform-workspace, time, user, venv, vgo, wsl)",
			"Unrecognized modules will be invoked as 'powerline-go-MODULE' executable plugins and should output a (possibly empty) list of JSON objects that unmarshal to powerline-go's Segment structs.")),
	Priority: flag.String(
		"priority",
		strings.Join(defaults.Priority, ","),
		commentsWithDefaults("Segments sorted by priority, if not enough space exists, the least priorized segments are removed first. Separate with ','",
			"(valid choices: aws, bzr, cwd, direnv, docker, docker-context, dotenv, duration, exit, fossil, gcp, git, gitlite, goenv, hg, host, jobs, kube, load, newline, nix-shell, node, perlbrew, perms, plenv, rbenv, root, rvm, shell-var, shenv, ssh, svn, termtitle, terraform-workspace, time, user, venv, vgo, vi-mode, wsl)")),
	MaxWidthPercentage: flag.Int(
		"max-width",
		defaults.MaxWidthPercentage,
		comments("Maximum width of the shell that the prompt may use, in percent. Setting this to 0 disables the shrinking subsystem.")),
	TruncateSegmentWidth: flag.Int(
		"truncate-segment-width",
		defaults.TruncateSegmentWidth,
		commentsWithDefaults("Maximum width of a segment, segments longer than this will be shortened if space is limited. Setting this to 0 disables it.")),
	PrevError: flag.Int(
		"error",
		defaults.PrevError,
		comments("Exit code of previously executed command")),
	NumericExitCodes: flag.Bool(
		"numeric-exit-codes",
		defaults.NumericExitCodes,
		comments("Shows numeric exit codes for errors.")),
	IgnoreRepos: flag.String(
		"ignore-repos",
		strings.Join(defaults.IgnoreRepos, ","),
		comments("A list of git repos to ignore. Separate with ','.",
			"Repos are identified by their root directory.")),
	ShortenGKENames: flag.Bool(
		"shorten-gke-names",
		defaults.ShortenGKENames,
		comments("Shortens names for GKE Kube clusters.")),
	ShortenEKSNames: flag.Bool(
		"shorten-eks-names",
		defaults.ShortenEKSNames,
		comments("Shortens names for EKS Kube clusters.")),
	ShortenOpenshiftNames: flag.Bool(
		"shorten-openshift-names",
		defaults.ShortenOpenshiftNames,
		comments("Shortens names for Openshift Kube clusters.")),
	ShellVar: flag.String(
		"shell-var",
		defaults.ShellVar,
		comments("A shell variable to add to the segments.")),
	ShellVarNoWarnEmpty: flag.Bool(
		"shell-var-no-warn-empty",
		defaults.ShellVarNoWarnEmpty,
		comments("Disables warning for empty shell variable.")),
	TrimADDomain: flag.Bool(
		"trim-ad-domain",
		defaults.TrimADDomain,
		comments("Trim the Domainname from the AD username.")),
	PathAliases: flag.String(
		"path-aliases",
		"",
		comments("One or more aliases from a path to a short name. Separate with ','.",
			"An alias maps a path like foo/bar/baz to a short name like FBB.",
			"Specify these as key/value pairs like foo/bar/baz=FBB.",
			"Use '~' for your home dir. You may need to escape this character to avoid shell substitution.")),
	Duration: flag.String(
		"duration",
		defaults.Duration,
		comments("The elapsed clock-time of the previous command")),
	Time: flag.String(
		"time",
		defaults.Time,
		comments("The layout string how a reference time should be represented.",
			"The reference time is predefined and not user choosen.",
			"Consult the golang documentation for details: https://pkg.go.dev/time#example-Time.Format")),
	DurationMin: flag.String(
		"duration-min",
		defaults.DurationMin,
		comments("The minimal time a command has to take before the duration segment is shown")),
	DurationLowPrecision: flag.Bool(
		"duration-low-precision",
		defaults.DurationLowPrecision,
		comments("Use low precision timing for duration with milliseconds as maximum resolution")),
	Eval: flag.Bool(
		"eval",
		defaults.Eval,
		comments("Output prompt in 'eval' format.")),
	Condensed: flag.Bool(
		"condensed",
		defaults.Condensed,
		comments("Remove spacing between segments")),
	IgnoreWarnings: flag.Bool(
		"ignore-warnings",
		defaults.IgnoreWarnings,
		comments("Ignores all warnings regarding unset or broken variables")),
	ViMode: flag.String(
		"vi-mode",
		defaults.ViMode,
		comments("The current vi-mode (eg. KEYMAP for zsh) for vi-module module")),
}
