package main

import (
	"html/template"
)

type Post struct {
	ID 			 int 		      `json:"id"`
	Content 	 template.HTML    `json:"content"`
	ReplyTo		 string			  `json:"reply_to"`
	ParentID 	 int 	     	  `json:"parent_id"`
	PublishedAt  string 		  `json:"published_at"`
	ChildrenIDs	 []int		  	  `json:"children_ids"`
}

type Reply struct {
	ReplyTo string `json:"reply_to"`
}
