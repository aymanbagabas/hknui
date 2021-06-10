package ui

import (
	"fmt"
	"strings"

	"github.com/muesli/reflow/ansi"
)

func (m model) statusBar() string {
	logo := statusBarLogoStyle.Render(" Y ") +
		statusBarNameStyle.Render("Hacker News ")

	tabs := ""
	tabNames := []string{"Top", "New", "Best", "Ask", "Show", "Jobs"}
	for i, tab := range tabNames {
		gap := tabStyle.Render(" ")
		sep := tabStyle.Render("|")
		underline := tabStyle.Bold(strings.ToLower(tab) == string(m.selectorType)).Underline(true).Render(tab[0:1])
		rest := tabStyle.Bold(strings.ToLower(tab) == string(m.selectorType)).Underline(false).Render(tab[1:])
		_tab := gap + underline + rest + gap
		if i != len(tabNames)-1 {
			_tab += sep
		}
		tabs += _tab
	}

	padding := max(0,
		m.viewport.Width-
			ansi.PrintableRuneWidth(logo)-
			ansi.PrintableRuneWidth(tabs),
	)
	emptySpace := tabStyle.Render(strings.Repeat(" ", padding))

	return fmt.Sprintf("%s%s%s", logo, tabs, emptySpace)
}
