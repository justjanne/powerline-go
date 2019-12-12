package powerline

import (
	runewidth "github.com/mattn/go-runewidth"
)

// Segment describes an information to display on the command line prompt
type Segment struct {
	// Content is the text to be displayed on the command line prompt
	Content string
	// Foreground is the text color (see https://misc.flogisoft.com/bash/tip_colors_and_formatting#background1)
	Foreground uint8
	// Background is the color of the filling background (see https://misc.flogisoft.com/bash/tip_colors_and_formatting#background1)
	Background uint8
	// Separator is the character to be used when generating multiple segments to override the default separator
	Separator string
	// SeparatorForeground is the character to be used when generating multiple segments to override the default foreground separator
	SeparatorForeground uint8
	// Priority is the priority of the segment. The higher, the less probable the segment will be dropped if the total length is too long
	Priority int
	// HideSeparators indicated not to display any separator with next segment.
	HideSeparators bool
	Width          int
}

func (s Segment) ComputeWidth(condensed bool) int {
	if condensed {
		return runewidth.StringWidth(s.Content) + runewidth.StringWidth(s.Separator)
	}
	return runewidth.StringWidth(s.Content) + runewidth.StringWidth(s.Separator) + 2
}
