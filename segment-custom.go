package main

import (
	"os"
	"strconv"
)

func segmentCustom(p *powerline) {
	customVarContent, customVarExists := os.LookupEnv("POWERLINE_CUSTOM_CONTENT")
	if customVarExists {
		customFG, _ := os.LookupEnv("POWERLINE_CUSTOM_FG")
		fg, _ := strconv.Atoi(customFG)
		customBG, _ := os.LookupEnv("POWERLINE_CUSTOM_BG")
		bg, _ := strconv.Atoi(customBG)

		p.appendSegment("custom", segment{
			content:    customVarContent,
			foreground: uint8(fg),
			background: uint8(bg),
		})
	}
}
