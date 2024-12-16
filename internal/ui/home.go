package ui

import (
	"v2ex-tui/internal/crawler"
	"v2ex-tui/internal/model"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomePage struct {
	topics   []model.Topic
	table    table.Model
	loading  bool
	err      error
	spinner  spinner.Model
	crawler  *crawler.Crawler
	selected int
}

func NewHomePage() *HomePage {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	columns := []table.Column{
		{Title: "📌 标题", Width: 100},
		{Title: "👤 作者", Width: 15},
		{Title: "💬 评论数", Width: 10},
		{Title: "🕒 时间", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	// 设置表格样式
	s1 := table.DefaultStyles()
	s1.Header = s1.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s1.Selected = s1.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	t.SetStyles(s1)

	return &HomePage{
		table:    t,
		loading:  true,
		spinner:  s,
		crawler:  crawler.New(),
		selected: 0,
	}
}

func (h *HomePage) Init() tea.Cmd {
	return tea.Batch(
		h.spinner.Tick,
		h.fetchTopics,
	)
}

func (h *HomePage) fetchTopics() tea.Msg {
	topics, err := h.crawler.FetchTopics()
	if err != nil {
		return errMsg{err}
	}
	return topicsMsg(topics)
}

type topicsMsg []model.Topic
type errMsg struct{ error }

func (h *HomePage) Update(msg tea.Msg) (*HomePage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "down":
			h.table, _ = h.table.Update(msg)
			h.selected = h.table.Cursor()
			return h, nil
		case "enter":
			if !h.loading && h.selected < len(h.topics) {
				return h, nil // 将在主程序中处理页面切换
			}
		case "r":
			h.loading = true
			return h, h.fetchTopics
		}

	case topicsMsg:
		h.loading = false
		h.topics = msg

		var rows []table.Row
		for _, t := range h.topics {
			rows = append(rows, table.Row{
				t.Title,
				t.Author,
				t.Comments,
				t.Time,
			})
		}
		h.table.SetRows(rows)
		return h, nil

	case errMsg:
		h.err = msg
		h.loading = false
		return h, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		h.spinner, cmd = h.spinner.Update(msg)
		return h, cmd
	}

	return h, nil
}

func (h *HomePage) View() string {
	if h.loading {
		return titleStyle.Render("V2EX 热门话题") + "\n" +
			h.spinner.View() + " 加载中...\n"
	}

	if h.err != nil {
		return errorStyle.Render("Error: "+h.err.Error()) + "\n"
	}

	return titleStyle.Render("V2EX 热门话题") + "\n" +
		tableStyle.Render(h.table.View()) + "\n\n" +
		subtitleStyle.Render(IconRefresh+" r 刷新 | "+IconEnter+" enter 查看详情 | q 退出\n")
}

func (h *HomePage) GetSelectedTopic() *model.Topic {
	if h.selected < len(h.topics) {
		return &h.topics[h.selected]
	}
	return nil
}
