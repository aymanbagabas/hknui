package ui

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/lukakerr/hkn"
)

func getItems(client hkn.Client, ids []int) (hkn.Items, error) {
	var items []hkn.Item
	for _, id := range ids {
		item, err := client.GetItem(id)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

func getStories(client hkn.Client, view selectorType, num int) ([]int, error) {
	var ids []int
	var err error

	switch view {
	case Top:
		ids, err = client.GetTopStories(num)
	case New:
		ids, err = client.GetNewStories(num)
	case Best:
		ids, err = client.GetBestStories(num)
	case Ask:
		ids, err = client.GetLatestAskStories(num)
	case Show:
		ids, err = client.GetLatestShowStories(num)
	case Job:
		ids, err = client.GetLatestJobStories(num)
	}

	if err != nil {
		log.Fatalf("Failed to fetch %s stories: %s\n", view, err.Error())
		return nil, err
	}

	return ids, nil
}
func (m model) updateItem(id int, force bool) tea.Cmd {
	return func() tea.Msg {
		var item *hkn.Item
		// log.Printf("getting item %d\n", id)
		if it, ok := m.items.Load(id); ok && !force {
			i := it.(hkn.Item)
			item = &i
		} else {
			it, err := m.client.GetItem(id)
			if err != nil {
				item = nil
			}
			item = &it
		}
		return itemMsg(item)
	}
}

func (m model) updateItems(ids []int, force bool) tea.Cmd {
	var cmds []tea.Cmd
	for _, id := range ids {
		cmds = append(cmds, m.updateItem(id, force))
	}
	return tea.Batch(cmds...)
}

func (m model) refreshStories(view selectorType) tea.Cmd {
	return func() tea.Msg {
		ids, _ := getStories(*m.client, view, fetchLimit)
		return storiesMsg(ids)
	}
}

func (m *model) getItemIdsInView() []int {
	if len(m.list) <= 0 {
		return []int{}
	}
	topPos := m.viewport.YOffset % (m.cursor + 1)
	botPos := topPos + m.viewport.Height
	if len(m.list) < botPos {
		botPos = len(m.list)
	}
	return m.list[topPos:botPos]
}

// func (m model) updateComments(ids []int) tea.Cmd {
// 	var cmds []tea.Cmd
// 	cmds = append(cmds, func() tea.Msg {
// 		kids := getItems(*m.client, ids)
// 		return itemsMsg{
// 			items:      kids,
// 			isComments: true,
// 		}
// 	})
// 	return tea.Batch(cmds...)
// }

func (m model) updateRenderer() tea.Cmd {
	return func() tea.Msg {
		r, _ := glamour.NewTermRenderer(
			glamour.WithStylePath("dark"), // FIXME style
			glamour.WithWordWrap(m.width),
		)
		return rendererMsg(r)
	}
}

func (m model) ensureItemsExist(ids []int) tea.Cmd {
	var cmds []tea.Cmd
	for _, id := range ids {
		// if item, ok := m.items[id]; ok {
		// 	// cmds = append(cmds, m.ensureItemsExist(item.Kids))
		// } else {
		cmds = append(cmds, m.updateItems([]int{id}, false))
		// }
	}
	return tea.Batch(cmds...)
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (m *model) stillLoading() bool {
	for _, id := range m.list {
		if _, ok := m.items.Load(id); !ok {
			return true
		}
	}
	return false
}

func (m model) numOfLines() int {
	return len(strings.Split(m.content, "\n"))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
