# A Powerline style prompt for your shell

A [Powerline](https://github.com/Lokaltog/vim-powerline) like prompt for Bash,
ZSH and Fish. Based on [Powerline-Shell](https://github.com/banga/powerline-shell) by @banga.
Ported to golang by @justjanne.

![MacVim+Solarized+Powerline+CtrlP](https://raw.github.com/banga/powerline-shell/master/bash-powerline-screenshot.png)

- Shows some important details about the git/hg branch (see below)
- Changes color if the last command exited with a failure code
- If you're too deep into a directory tree, shortens the displayed path with an ellipsis
- Shows the current Python [virtualenv](http://www.virtualenv.org/) environment
- It's easy to customize and extend. See below for details.

**Table of Contents** 

- [Version Control](#version-control)
- [Setup](#setup)
  - [All Shells](#all-shells)
  - [Bash](#bash)
  - [ZSH](#zsh)
  - [Fish](#fish)

## Version Control

All of the version control systems supported by powerline shell give you a
quick look into the state of your repo:

- The current branch is displayed and changes background color when the
  branch is dirty.
- When the local branch differs from the remote, the difference in number
  of commits is shown along with `⇡` or `⇣` indicating whether a git push
  or pull is pending

In addition, git has a few extra symbols:

- `✎` -- a file has been modified, but not staged for commit
- `✔` -- a file is staged for commit
- `✼` -- a file has conflicts
- `+` -- untracked files are present

Each of these will have a number next to it if more than one file matches.

## Setup

This script uses ANSI color codes to display colors in a terminal. These are
notoriously non-portable, so may not work for you out of the box, but try
setting your $TERM to `xterm-256color`, because that works for me.

- Patch the font you use for your terminal: see
  [powerline-fonts](https://github.com/Lokaltog/powerline-fonts)
  - If you struggle too much to get working fonts in your terminal, you can use
    "compatible" mode.
  - If you're using old patched fonts, you have to use the older symbols.
    Basically reverse [this
    commit](https://github.com/milkbikis/powerline-shell/commit/2a84ecc) in
    your copy

- Set your GOPATH (otherwise the package will be installed in $HOME/go/bin/)

- Download and install the package

```
go get github.com/justjanne/powerline-go
```

- Move it to a convenient place (for example, `~/.powerline/powerline-go`)

### All Shells

There are a few optional arguments which can be seen by running
`powerline-go -help`.

```
Usage of powerline-go:
  -colorize-hostname
        Colorize the hostname based on a hash of itself
  -cwd-max-depth int
        Maximum number of directories to show in path
        (default 5)
  -cwd-max-dir-size int
        Maximum number of letters displayed for each directory in the path
        (default -1)
  -cwd-mode string
        How to display the current directory
        (default "fancy")
  -error int
        Exit code of previously executed command
  -mode string
        The characters used to make separators between segments.
        (valid choices: patched, compatible, flat)
        (default "patched")
  -modules string
        The list of modules to load. Separate with ','
        (valid choices: cwd, exit, git, hg, host, jobs, perms, root, ssh, user, venv)
        (default "venv,user,ssh,host,cwd,perms,jobs,exit,root")
  -shell string
        Set this to your shell type
        (valid choices: bare, bash, zsh)
        (default "bash")
```

### Bash

Add the following to your `.bashrc` (or `.profile` on Mac):

```
function _update_ps1() {
    PS1="$(~/.powerline/powerline-go -error $?)"
}

if [ "$TERM" != "linux" ]; then
    PROMPT_COMMAND="_update_ps1; $PROMPT_COMMAND"
fi
```

### ZSH

Add the following to your `.zshrc`:

```
function powerline_precmd() {
    PS1="$(~/.powerline/powerline-go -error $? -shell zsh)"
}

function install_powerline_precmd() {
  for s in "${precmd_functions[@]}"; do
    if [ "$s" = "powerline_precmd" ]; then
      return
    fi
  done
  precmd_functions+=(powerline_precmd)
}

if [ "$TERM" != "linux" ]; then
    install_powerline_precmd
fi
```

### Fish

Redefine `fish_prompt` in `~/.config/fish/config.fish:`

```
function fish_prompt
    ~/.powerline/powerline-go -error $status -shell bare
end
```

## Customization

### Adding, Removing and Re-arranging segments

You can use the `-modules` argument to define which segments to draw, and in
which order. Segments can also be repeated. Simply change the argument in your
shell’s init to customize the prompt.