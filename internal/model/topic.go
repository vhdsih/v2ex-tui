package model

// Topic represents a V2EX topic
type Topic struct {
	Title    string
	Author   string
	Comments string
	Time     string
	Content  string
	URL      string
	Replies  []Reply
}
