# A Powerline style prompt for your shell

A [Powerline](https://github.com/Lokaltog/vim-powerline) like prompt for Bash,
ZSH and Fish. Based on [Powerline-Shell](https://github.com/banga/powerline-shell) by @banga.
Ported to golang by @justjanne.

![Solarized+Powerline](https://raw.github.com/justjanne/powerline-go/master/preview.png)

- Shows some important details about the git/hg branch (see below)
- Changes color if the last command exited with a failure code
- If you're too deep into a directory tree, shortens the displayed path with an ellipsis
- Shows the current Python [virtualenv](http://www.virtualenv.org/) environment
- It's easy to customize and extend. See below for details.

**Table of Contents** 

- [Version Control](#version-control)
- [Installation](#installation)
  - [Precompiled Binaries](#precompiled-binaries)
  - [Other Platforms](#other-platforms)
  - [Bash](#bash)
  - [ZSH](#zsh)
  - [Fish](#fish)
- [Customization](#customization)
- [License](#license)

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

## Installation

`powerline-go` uses ANSI color codes, these should nowadays work everywhere,
but you may have to set your $TERM to `xterm-256color` for it to work.

If you want to use the "patched" mode (which is the default, and provides
improved UI), you'll need to install a powerline font, either as fallback,
or by patching the font you use for your terminal: see
[powerline-fonts](https://github.com/Lokaltog/powerline-fonts).  
Alternatively you can use "compatible" or "flat" mode.

### Precompiled Binaries

I provide precompiled binaries for x64 Linux and macOS in the
[releases tab](https://github.com/justjanne/powerline-go/releases)

### Other Platforms

- Install (and update) the package with

```
go get -u github.com/justjanne/powerline-go
```

- By default it will be in `~/go/bin`, if you want to change that, you can set
  your `$GOPATH` and/or `$GOBIN`, but will need ot change the path in the
  following scripts, too.

### Bash

Add the following to your `.bashrc` (or `.profile` on Mac):

```
function _update_ps1() {
    if [ -e $GOPATH/bin/powerline-go ]; then
        PS1="$($GOPATH/bin/powerline-go -error $?)"
    fi
}

if [ "$TERM" != "linux" ]; then
    PROMPT_COMMAND="_update_ps1; $PROMPT_COMMAND"
fi
```

### ZSH

Add the following to your `.zshrc`:

```
function powerline_precmd() {
    PS1="$(~/go/bin/powerline-go -error $? -shell zsh)"
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
    ~/go/bin/powerline-go -error $status -shell bare
end
```

## Customization

There are a few optional arguments which can be seen by running
`powerline-go -help`. These can be used by changing the command you have set
in your shell’s init file.

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
  -east-asian-width
        Use East Asian Ambiguous Widths
  -error int
        Exit code of previously executed command
  -max-width int
        Maximum width of the shell that the prompt may use, in percent. Setting this to 0 disables the shrinking subsystem.
        (default 50)
  -mode string
        The characters used to make separators between segments.
        (valid choices: patched, compatible, flat)
        (default "patched")
  -modules string
        The list of modules to load. Separate with ','
        (valid choices: cwd, docker, exit, git, hg, host, jobs, perms, root, ssh, time, user, venv)
        (default "venv,user,host,ssh,cwd,perms,git,hg,jobs,exit,root")
  -newline
        Show the prompt on a new line
  -priority string
        Segments sorted by priority, if not enough space exists, the least priorized segments are removed first. Separate with ','
        (valid choices: cwd, docker, exit, git-branch, git-status, hg, host, jobs, perms, root, ssh, time, user, venv)
        (default "root,cwd,user,host,ssh,perms,git-branch,git-status,hg,jobs,exit")
  -shell string
        Set this to your shell type
        (valid choices: bare, bash, zsh)
        (default "bash")
  -theme string
        Set this to the theme you want to use
        (valid choices: default)
        (default "default")
```


## License

> This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.  
> This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.  
> You should have received a copy of the GNU General Public License along with this program. If not, see <http://www.gnu.org/licenses/>.  
