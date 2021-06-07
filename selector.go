package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/lukakerr/hkn"
	"github.com/muesli/reflow/wordwrap"
)

type selectorType string

const (
	Top  selectorType = "top"
	New  selectorType = "new"
	Best selectorType = "best"
	Ask  selectorType = "ask"
	Show selectorType = "show"
	Job  selectorType = "job"
)

func selectFunc(item hkn.Item, index int, selected bool) string {
	num := selectorNumStyle.Render(fmt.Sprintf("%d. ", index+1))
	titleStyle := lg.NewStyle()
	if selected {
		titleStyle = selectorTitleActiveStyle
	} else {
		titleStyle = selectorTitleStyle
	}
	title := titleStyle.Render(item.Title)
	baseUrl := ""
	if item.URL != "" {
		url, err := url.Parse(item.URL)
		if err == nil {
			baseUrl = selectorUrlStyle.Render(" (" + url.Hostname() + ")")
		}
	}

	return fmt.Sprintf("%s%s%s\n", num, title, baseUrl)
}

func (m model) getItem(id int) *hkn.Item {
	val, ok := m.items.Load(id)
	if ok {
		item := val
		return item.(*hkn.Item)
	} else {
		return nil
	}
}

func (m model) getCursorItemId() int {
	var id int = -1
	if len(m.list) > 0 {
		id = m.list[m.cursor]
	}
	return id
}

func (m model) getItemUnderCursor() *hkn.Item {
	item := m.getItem(m.getCursorItemId())
	return item
}

func (m *model) updateFooter() string {
	item, ok := m.items.Load(m.itemId)
	if ok {
		item := item.(hkn.Item)
		time := time.Unix(int64(item.Time), 0).Local()

		rv := wordwrap.String(
			fmt.Sprintf("%s points by %s %s | %d comment%s",
				lg.NewStyle().Foreground(lg.Color("15")).Render(fmt.Sprintf("%d", item.Score)),
				lg.NewStyle().Foreground(lg.Color("3")).Render(item.By),
				timeString(time),
				item.Descendants,
				map[bool]string{true: "s", false: ""}[len(item.Kids) != 1],
			),
			m.width,
		)
		// log.Print("scroll percent")
		// log.Print(len(strings.Split(m.viewport.View(), "\n")))
		// log.Printf("%f %d %d %d\n", m.viewport.ScrollPercent(), m.viewport.Height, m.viewport.Width, m.viewport.YOffset)
		scrollPercent := fmt.Sprintf("%d%%", int(m.viewport.ScrollPercent()*100))
		repeatWidth := m.width - lg.Width(rv) - lg.Width(scrollPercent)
		if repeatWidth < 0 {
			scrollPercent = "0%"
		}
		padding := strings.Repeat(" ", repeatWidth)
		footer := fmt.Sprintf("%s%s%s", rv, padding, scrollPercent)
		m.setFooter(footer)

		return footer
	}
	return ""
}

func (m model) selectorView() string {
	var s string
	for i, id := range m.list {
		item, ok := m.items.Load(id)
		if ok {
			s += selectFunc(item.(hkn.Item), i, i == m.cursor)
		}
	}

	return s
}
