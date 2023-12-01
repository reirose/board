package src

import (
	"html/template"
)

type Post struct {
	ID          int           `json:"id"`
	Content     template.HTML `json:"content"`
	ReplyTo     string        `json:"-"`
	ParentID    int           `json:"-"`
	Children    []*PostReply  `json:"-"`
	PublishedAt string        `json:"published_at"`
	ChildrenIDs []int         `json:"children_ids"`
}

type PostReply struct {
	ID          int           `json:"id"`
	Content     template.HTML `json:"content"`
	ChildrenIDs []int         `json:"children_ids"`
	PublishedAt string        `json:"published_at"`
}

type Reply struct {
	ReplyTo string `json:"reply_to"`
}

type User struct {
	ID       int    `json:"id"`
	Role     string `json:"role"`
	UserID   string `json:"user_id"`
	Password string `json:"-"`
	Token    string `json:"-"`
}

type PreRegUser struct {
	Role   string `json:"role"`
	UserID string `json:"user_id"`
}

type APIResponse struct {
	Version string  `json:"version"`
	Users   []*User `json:"users"`
	Posts   []*Post `json:"posts"`
	JSON    []byte  `json:"-"`
}
