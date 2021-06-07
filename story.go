package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/lukakerr/hkn"
)

var converter = md.NewConverter("", true, &md.Options{
	EmDelimiter:     "*",
	StrongDelimiter: "**",
})

func (m model) storyView() string {
	item, ok := m.items.Load(m.itemId)
	if ok {
		item := item.(hkn.Item)
		title := storyTitleStyle.Render(item.Title)
		var text string
		if item.Text != "" {
			text = item.Text
		} else {
			text = item.URL
		}
		text = m.storyRender(text)

		sep := strings.Repeat("─", m.width)

		return fmt.Sprintf("%s\n\n%s\n%s\n\n%s", title, text, sep, m.buildComments(item, 0))
	}
	return ""
}

func (m *model) buildComments(item hkn.Item, depadding int) string {
	var s strings.Builder
	for _, id := range item.Kids {
		val, ok := m.items.Load(id)
		if ok {
			kid := val.(hkn.Item)
			user := storyCommentUserStyle.Render(kid.By)
			actualTime := time.Unix(int64(kid.Time), 0).Local()
			postingTime := timeString(actualTime)
			fmt.Fprintf(&s, "%s %s\n", user, storyCommentTimeStyle.Render(postingTime))
			fmt.Fprintf(&s, "%s\n", m.storyRender(kid.Text))
			subComments := m.buildComments(kid, depadding+2)
			if subComments != "" {
				fmt.Fprintf(&s, "%s\n", storyCommentKidsStyle.Width(m.viewport.Width-depadding).Render(subComments))
			}
		} else {
			log.Print(fmt.Sprintf("%d not found\n", id))
		}
	}

	return s.String()
}

func (m *model) storyRender(text string) string {
	markdown, _ := m.md.Render(renderText(text))
	return markdown
}

func timeString(actualTime time.Time) string {
	since := time.Since(actualTime)
	seconds := int(since.Seconds())
	minutes := int(since.Minutes())
	hours := int(since.Hours())
	days := int(hours / 24)
	months := int(days / 30)
	postingTime := fmt.Sprintf("on %s %d, %d", actualTime.Month().String(), actualTime.Day(), actualTime.Year())
	plural := map[bool]string{true: "s", false: ""}
	if seconds < 60 {
		postingTime = fmt.Sprintf("%d second%s ago", seconds, plural[seconds != 1])
	} else if minutes < 60 {
		postingTime = fmt.Sprintf("%d minute%s ago", minutes, plural[minutes != 1])
	} else if hours < 24 {
		postingTime = fmt.Sprintf("%d hour%s ago", hours, plural[hours != 1])
	} else if days < 30 {
		postingTime = fmt.Sprintf("%d day%s ago", days, plural[days != 1])
	} else if months < 12 {
		postingTime = fmt.Sprintf("%d month%s", months, plural[months != 1])
	}
	return postingTime
}

func renderText(text string) string {
	// text = html.UnescapeString(h2t.HTML2Text(s))
	// return text
	markdown, err := converter.ConvertString(text)
	if err != nil {
		return text
	}
	return markdown
}
