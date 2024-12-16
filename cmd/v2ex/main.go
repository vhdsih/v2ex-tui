package main

import (
	"fmt"

	"v2ex-tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

type page int

const (
	homeView page = iota
	detailView
	replyView
)

type model struct {
	currentPage page
	homePage    *ui.HomePage
	detailPage  *ui.DetailPage
	replyPage   *ui.ReplyPage
}

func initialModel() model {
	return model{
		currentPage: homeView,
		homePage:    ui.NewHomePage(),
		detailPage:  ui.NewDetailPage(),
		replyPage:   ui.NewReplyPage(),
	}
}

func (m model) Init() tea.Cmd {
	return m.homePage.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.currentPage == replyView {
				m.currentPage = detailView
				return m, nil
			}
			if m.currentPage == detailView {
				m.currentPage = homeView
				return m, nil
			}
		case "enter":
			if m.currentPage == homeView {
				if topic := m.homePage.GetSelectedTopic(); topic != nil {
					m.currentPage = detailView
					return m, m.detailPage.LoadTopic(*topic)
				}
			} else if m.currentPage == detailView {
				if reply := m.detailPage.GetSelectedReply(); reply != nil {
					m.currentPage = replyView
					m.replyPage.LoadReply(*reply, m.detailPage.Topic.Replies)
					return m, nil
				}
			}
		}
	}

	var cmd tea.Cmd
	switch m.currentPage {
	case homeView:
		m.homePage, cmd = m.homePage.Update(msg)
	case detailView:
		m.detailPage, cmd = m.detailPage.Update(msg)
	case replyView:
		m.replyPage, cmd = m.replyPage.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	switch m.currentPage {
	case homeView:
		return m.homePage.View()
	case detailView:
		return m.detailPage.View()
	case replyView:
		return m.replyPage.View()
	default:
		return "Unknown view"
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		return
	}
}
