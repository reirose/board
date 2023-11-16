package main

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	log(r)
	posts, err := dbGetAllPosts()
	catch(err)

	t, _ := template.ParseFiles("templates/base.html", "templates/index.html")
	err = t.Execute(w, posts)
	catch(err)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	log(r)
	post := r.Context().Value("post").(*Post)

	t, _ := template.ParseFiles("templates/base.html", "templates/post.html")
	err := t.Execute(w, post)
	if err != nil {
		print(err)
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	log(r)
	replyTo := getParam(r.URL.String(), "reply_to")
	reply := &Reply{ReplyTo: replyTo}

	t, _ := template.ParseFiles("templates/base.html", "templates/create.html")
	err := t.Execute(w, reply)
	catch(err)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	log(r)
	t := time.Now()
	content := r.FormValue("content")
	publishedAt := t.UTC().String()
	parentId := r.FormValue("parent_id")

	parentIdInt, err := strconv.Atoi(parentId)
	if err != nil {
		parentIdInt = 0
	}

	post := &Post{
		Content:     template.HTML(content),
		PublishedAt: publishedAt,
		ParentID:    parentIdInt,
	}

	err = dbCreatePost(post)
	catch(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

func ErrorNotFound(w http.ResponseWriter, r *http.Request) {
	log(r)

	t, _ := template.ParseFiles("templates/base.html", "templates/notfound.html")
	err := t.Execute(w, nil)
	catch(err)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	log(r)
	post := r.Context().Value("post").(*Post)
	err := dbDeletePost(post.ID)
	catch(err)

	http.Redirect(w, r, "/", http.StatusFound)
}
