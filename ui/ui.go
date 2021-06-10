package ui

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/lukakerr/hkn"
	"github.com/pkg/browser"
)

const (
	fetchLimit = 200
)

type view string

const (
	ArticleView view = "article"
	StoryView   view = "story"
)

const (
	headerHeight               = 2
	footerHeight               = 2
	useHighPerformanceRenderer = false
)

type itemMsg *hkn.Item

type storiesMsg []int

type rendererMsg *glamour.TermRenderer

type model struct {
	loading bool
	ready   bool
	client  *hkn.Client
	width   int
	height  int
	view    view
	spinner spinner.Model

	itemId int
	md     *glamour.TermRenderer

	items        *sync.Map
	cursor       int
	list         []int
	selectorType selectorType
	prevOffset   int

	viewport *viewport.Model

	header  string
	content string
	footer  string
}

func NewModel() tea.Model {
	client := hkn.NewClient()
	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

	return &model{
		loading:      false,
		ready:        false,
		client:       client,
		view:         StoryView,
		items:        &sync.Map{},
		cursor:       0,
		list:         make([]int, 0),
		selectorType: Top,
		spinner:      s,
		viewport:     &viewport.Model{Width: 0, Height: 0, HighPerformanceRendering: useHighPerformanceRenderer},
	}
}

func (m *model) initViewport(width, height int) {
	margins := headerHeight + footerHeight
	m.viewport = &viewport.Model{Width: width, Height: height - margins}
	m.viewport.YPosition = headerHeight + 1
	m.viewport.YOffset = 0
	m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
	m.viewport.SetContent(m.content)
}

func (m *model) resizeViewport(width, height int) {
	margins := headerHeight + footerHeight
	m.viewport.Width = width
	m.viewport.Height = height - margins
}

func (m *model) setFooter(footer string) {
	m.footer = footer
}

func (m *model) viewUp(n int) {
	switch m.view {
	case StoryView:
		if m.cursor > 0 {
			m.cursor -= n
			m.viewport.LineUp(n)
		}
		if m.cursor < 0 {
			m.cursor = 0
			m.viewport.GotoTop()
		}
	case ArticleView:
		m.viewport.LineUp(n)
	}
}

func (m *model) viewDown(n int) {
	if item, ok := m.items.Load(m.itemId); ok {
		log.Print(item.(hkn.Item))
	}
	switch m.view {
	case StoryView:
		lines := m.numOfLines()
		if m.cursor+n < len(m.list) {
			m.cursor += n
			m.viewport.LineDown(n)
		}
		if m.cursor+n >= lines {
			m.cursor = lines - 2
			m.viewport.GotoBottom()
		}
	case ArticleView:
		m.viewport.LineDown(n)
	}
}

func (m *model) updateItemUnderCursor(force bool) tea.Cmd {
	var cmds []tea.Cmd
	if m.cursor >= 0 && m.cursor < len(m.list) {
		itemId := m.list[m.cursor]
		item, ok := m.items.Load(itemId)
		if ok {
			item := item.(hkn.Item)
			m.itemId = itemId
			cmds = append(cmds, m.updateItem(itemId, force))
			cmds = append(cmds, m.updateItems(item.Kids, force))
		}
	}
	return tea.Batch(cmds...)
}

func (m *model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		func() tea.Msg {
			return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("T")}
		},
	}
	return tea.Batch(cmds...)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updateViewport := false
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case rendererMsg:
		m.md = msg
	case storiesMsg:
		m.list = msg
		cmds = append(cmds, m.updateItems(m.list, false))
		cmds = append(cmds, m.updateItemUnderCursor(true))
	case itemMsg:
		if msg == nil {
			break
		}
		item := *msg
		m.items.Store(item.ID, item)
		switch item.Type {
		case "story":
			if m.view == StoryView {
				updateViewport = true
			}
			cmds = append(cmds, m.updateItems(item.Kids, false))
		case "comment":
			cmds = append(cmds, m.updateItems(item.Kids, false))
		}

	case tea.MouseMsg:
		// switch msg.String() {
		// case tea.MouseWheelUp:
		// 	m.selectorScrollUp(1)

		// case tea.MouseWheelDown:
		// 	m.selectorScrollDown(1)
		// }

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

		switch m.view {
		case ArticleView:
			switch msg.String() {
			case "left", "esc", "h":
				m.view = StoryView
				m.viewport.YOffset = m.prevOffset
				m.prevOffset = 0

			case "down", "j":
				m.viewDown(1)

			case "up", "k":
				m.viewUp(1)

			case "pgup", "u":
				m.viewUp(m.viewport.Height / 2)

			case "pgdown", "d":
				m.viewDown(m.viewport.Height / 2)

			case "g":
				m.viewport.GotoTop()

			case "G":
				m.viewport.GotoBottom()

			}
			cmds = append(cmds, m.updateItemUnderCursor(false))

		case StoryView:
			switch msg.String() {
			case "T", "N", "B", "A", "S", "J":
				switch msg.String() {
				case "T":
					m.selectorType = Top
				case "N":
					m.selectorType = New
				case "B":
					m.selectorType = Best
				case "A":
					m.selectorType = Ask
				case "S":
					m.selectorType = Show
				case "J":
					m.selectorType = Job
				}
				m.cursor = 0
				m.viewport.GotoTop()
				cmds = append(cmds, m.refreshStories(m.selectorType))
				cmds = append(cmds, m.updateItemUnderCursor(true))
				// m.loading = true

			case "up", "k":
				m.viewUp(1)
				cmds = append(cmds, m.updateItemUnderCursor(true))

			case "down", "j":
				m.viewDown(1)
				cmds = append(cmds, m.updateItemUnderCursor(true))

			case "pgup", "u":
				m.viewUp(m.viewport.Height / 2)
				cmds = append(cmds, m.updateItemUnderCursor(true))

			case "pgdown", "d":
				m.viewDown(m.viewport.Height / 2)
				cmds = append(cmds, m.updateItemUnderCursor(true))

			case "g":
				m.cursor = 0
				m.viewport.GotoTop()

			case "G":
				m.cursor = len(m.list) - 1
				m.viewport.GotoBottom()

			case "right", "l":
				m.view = ArticleView
				m.prevOffset = m.viewport.YOffset
				m.viewport.GotoTop()

			case "enter":
				item, ok := m.items.Load(m.itemId)
				if ok {
					item := item.(hkn.Item)
					if item.Text == "" { // URL
						browser.OpenURL(item.URL)
					} else {
						m.view = ArticleView
						m.prevOffset = m.viewport.YOffset
						m.viewport.GotoTop()
					}
				}

			case "o":
				browser.OpenURL(fmt.Sprintf("https://news.ycombinator.com/item?id=%d", m.itemId))

			}
		}
		updateViewport = true

	case tea.WindowSizeMsg:
		width := msg.Width
		height := msg.Height
		m.width = width
		m.height = height
		if !m.ready {
			m.initViewport(width, height)
			m.ready = true
		} else {
			m.resizeViewport(width, height)
		}
		cmds = append(cmds, m.updateRenderer())
		// Update viewport on init or resize
		updateViewport = true
	}

	if m.cursor >= 0 && m.cursor < len(m.list) {
		itemId := m.list[m.cursor]
		_, ok := m.items.Load(itemId)
		if ok {
			m.itemId = itemId
		}
	}

	if updateViewport {
		m.header = m.statusBar()
		switch m.view {
		case StoryView:
			// cmds = append(cmds, m.ensureItemsExist(m.list))
			m.content = m.selectorView()

		case ArticleView:
			m.content = m.storyView()
		}
		// if m.stillLoading() {
		// 	m.content = m.spinner.View()
		// 	m.spinner.Update(msg)
		// }
		m.viewport.SetContent(m.content)
		m.updateFooter()

		if m.viewport.HighPerformanceRendering {
			cmds = append(cmds, viewport.Sync(*m.viewport))
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	header := fmt.Sprintf("%s\n%s", m.header, strings.Repeat(" ", m.width))
	footer := fmt.Sprintf("\n%s", m.footer)

	return fmt.Sprintf("%s\n%s\n%s", header, m.viewport.View(), footer)
}
