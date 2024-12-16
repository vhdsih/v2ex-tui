package model

// Reply represents a reply to a topic
type Reply struct {
	Author  string
	Time    string
	Content string
	Number  string
	// 新增字段，用于追踪被回复的评论
	ReplyTo    string
	ReplyCount int
}
