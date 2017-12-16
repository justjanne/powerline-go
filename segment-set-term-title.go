package main

func segmentSetTermTitle(p *powerline) {
	var setTitle string

	if *p.args.Shell == "bash" {
		setTitle = "\\[\\e]0;\\u@\\h: \\w\\a\\]"
	} else if *p.args.Shell == "zsh" {
		setTitle = "%{\033]0;%n@%m: %~\007%}"
	}

	p.appendSegment("set-term-title", segment{
		content:  setTitle,
		priority: MaxInteger, // do not truncate
		special:  true,       // do not draw separators
	})
}
