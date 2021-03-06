package main

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
)

const (
	// LeadingSpaceNum is how many spaches before output
	LeadingSpaceNum = 2
)

var (
	// WhiteBoldString predefined color
	WhiteBoldString = color.New(color.FgWhite, color.Bold).SprintFunc()
	// WhiteString predefined color
	WhiteString = color.WhiteString
	// GreenString predefined color
	GreenString = color.GreenString
	// RedString predefined color
	RedString = color.RedString
	// GrayString predefined color
	GrayString = color.New(color.FgWhite, color.Faint).SprintfFunc()
	// LeadingSpace predefined string
	LeadingSpace = strings.Repeat(" ", LeadingSpaceNum)
	// ErrorOutput predefined string
	ErrorOutput = color.New(color.FgRed)
)

func output(page []string) (rendered string) {
	re1 := regexp.MustCompile("{{")
	re2 := regexp.MustCompile("}}")

	rendered += "\n"
	for _, line := range page {
		switch {
		case len(line) == 0:
			continue
		case line[0] == '#':
			rendered += LeadingSpace + WhiteBoldString(strings.ReplaceAll(line, "#", "")) + "\n\n"
		case line[0] == '>':
			rendered += LeadingSpace + WhiteString(strings.ReplaceAll(strings.ReplaceAll(line, ">", ""), "<", "")) + "\n"
		case line[0] == '-':
			rendered += "\n" + LeadingSpace + GreenString(line) + "\n"
		case line[0] == '`':
			line = line[1 : len(line)-1]
			res := re1.ReplaceAllString(line, "\n{{")
			res = re2.ReplaceAllString(res, "}}\n")
			rendered += LeadingSpace + LeadingSpace
			for _, item := range strings.Split(res, "\n") {
				switch {
				case len(item) == 0:
					rendered += " "
				case item[0] == '{':
					// If argument, print without color
					rendered += item[2 : len(item)-2]
				default:
					rendered += RedString(item)
				}
			}
			rendered += "\n"
		default:
			// new to add unmatched lines to output
			rendered += line
		}
	}
	rendered += "\n"
	return
}
