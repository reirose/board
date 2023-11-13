package main

import (
	"html/template"
	"net/http"
	"net/url"
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
	catch(err)
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	log(r)
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	params := u.Query()
	reply_to := params.Get("reply_to")
	reply := &Reply{ReplyTo: reply_to,}

	t, _ := template.ParseFiles("templates/base.html", "templates/create.html")
	err = t.Execute(w, reply)
	catch(err)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	log(r)
	t := time.Now()
	content := r.FormValue("content")
	published_at := t.UTC().String()
	parent_id := r.FormValue("parent_id")

	parent_id_int, err := strconv.Atoi(parent_id)
	if err != nil {
		parent_id_int = 0
	}

	post := &Post{
		Content:     template.HTML(content),
		PublishedAt: published_at,
		ParentID:    parent_id_int,
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
