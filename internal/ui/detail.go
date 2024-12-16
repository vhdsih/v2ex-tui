package ui

import (
	"fmt"
	"strings"

	"v2ex-tui/internal/crawler"
	"v2ex-tui/internal/model"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DetailPage struct {
	Topic    model.Topic
	table    table.Model
	loading  bool
	err      error
	spinner  spinner.Model
	crawler  *crawler.Crawler
	selected int
}

func NewDetailPage() *DetailPage {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	columns := []table.Column{
		{Title: "序号", Width: 10},
		{Title: "作者", Width: 15},
		{Title: "内容", Width: 80},
		{Title: "时间", Width: 20},
		{Title: "回复数", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

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

	return &DetailPage{
		table:    t,
		loading:  true,
		spinner:  s,
		crawler:  crawler.New(),
		selected: 0,
	}
}

func (d *DetailPage) LoadTopic(topic model.Topic) tea.Cmd {
	d.Topic = topic
	d.loading = true
	return d.fetchTopicDetail
}

func (d *DetailPage) fetchTopicDetail() tea.Msg {
	topic, err := d.crawler.FetchTopicDetail(d.Topic.URL)
	if err != nil {
		return errMsg{err}
	}

	// 统计每条评论的回复数
	for i := range topic.Replies {
		replyCount := 0
		currentAuthor := topic.Replies[i].Author
		// 遍历后续的评论来统计回复数
		for _, reply := range topic.Replies[i+1:] {
			if strings.HasPrefix(strings.TrimSpace(reply.Content), "@"+currentAuthor) {
				replyCount++
			}
		}
		topic.Replies[i].ReplyCount = replyCount
	}

	return topicDetailMsg{*topic}
}

type topicDetailMsg struct {
	topic model.Topic
}

func (d *DetailPage) Update(msg tea.Msg) (*DetailPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "down":
			d.table, _ = d.table.Update(msg)
			d.selected = d.table.Cursor()
			return d, nil
		case "enter":
			if !d.loading && len(d.Topic.Replies) > 0 {
				return d, nil // 将在主程序中处理页面切换
			}
		}

	case topicDetailMsg:
		d.loading = false
		d.Topic = msg.topic

		var rows []table.Row
		for _, r := range d.Topic.Replies {
			rows = append(rows, table.Row{
				r.Number,
				r.Author,
				r.Content,
				r.Time,
				fmt.Sprintf("%d", r.ReplyCount),
			})
		}
		d.table.SetRows(rows)
		return d, nil

	case errMsg:
		d.err = msg
		d.loading = false
		return d, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		d.spinner, cmd = d.spinner.Update(msg)
		return d, cmd
	}

	return d, nil
}

func (d *DetailPage) View() string {
	if d.loading {
		return titleStyle.Render("话题详情") + "\n" +
			d.spinner.View() + " 加载中...\n"
	}

	if d.err != nil {
		return errorStyle.Render("Error: "+d.err.Error()) + "\n"
	}

	var s string
	// 标题区域
	s += sectionStyle.Render(
		titleStyle.Render(IconTitle+"标题: "+d.Topic.Title)+"\n"+
			subtitleStyle.Render(IconAuthor+"作者: "+d.Topic.Author)+"\n"+
			subtitleStyle.Render(IconTime+"时间: "+d.Topic.Time),
	) + "\n"

	// 内容区域
	s += sectionStyle.Render(
		titleStyle.Render(IconContent+"内容:")+"\n"+
			contentStyle.Render(d.Topic.Content),
	) + "\n"

	// 评论区域
	s += titleStyle.Render(IconComments+" 评论:") + "\n" +
		tableStyle.Render(d.table.View()) + "\n\n" +
		subtitleStyle.Render(IconBack+" esc 返回 | "+IconEnter+" enter 查看评论详情 | q 退出\n")

	return s
}

func (d *DetailPage) GetSelectedReply() *model.Reply {
	if d.selected < len(d.Topic.Replies) {
		return &d.Topic.Replies[d.selected]
	}
	return nil
}
