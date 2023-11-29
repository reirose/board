package post

import (
	"board/src"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	posts, err := DbGetAllPosts()
	src.Catch(err)

	t, _ := template.ParseFiles("templates/base.html", "templates/index.html")
	err = t.Execute(w, posts)
	src.Catch(err)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	post := r.Context().Value("post").(*src.Post)

	t, _ := template.ParseFiles("templates/base.html", "templates/post.html")
	err := t.Execute(w, post)
	if err != nil {
		print(err)
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	replyTo := src.GetParam(r.URL.String(), "reply_to")
	reply := &src.Reply{ReplyTo: replyTo}

	t, _ := template.ParseFiles("templates/base.html", "templates/create.html")
	err := t.Execute(w, reply)
	src.Catch(err)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	t := time.Now()
	content := r.FormValue("content")
	publishedAt := t.UTC().String()
	parentId := r.FormValue("parent_id")

	parentIdInt, err := strconv.Atoi(parentId)
	if err != nil {
		parentIdInt = 0
	}

	post := &src.Post{
		Content:     template.HTML(content),
		PublishedAt: publishedAt,
		ParentID:    parentIdInt,
	}

	err = DbCreatePost(post)
	src.Catch(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

func ErrorNotFound(w http.ResponseWriter, r *http.Request) {
	src.Log(r)

	t, _ := template.ParseFiles("templates/base.html", "templates/notfound.html")
	err := t.Execute(w, nil)
	src.Catch(err)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	post := r.Context().Value("post").(*src.Post)
	err := DbDeletePost(post.ID)
	src.Catch(err)

	http.Redirect(w, r, "/", http.StatusFound)
}
