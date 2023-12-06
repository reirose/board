package post

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"board/src"
)

func GetRequestDataC(t []string, r *http.Request) *src.ReqData {
	src.Log(r)
	reqData := new(src.ReqData)
	for _, v := range t {
		switch v {
		case "user":
			user := new(src.User)
			userCtx := r.Context().Value("user")
			if userCtx != nil {
				user = userCtx.(*src.User)
			}

			reqData.UserData = user
		case "post":
			post := new(src.Post)
			postCtx := r.Context().Value("post")
			if postCtx != nil {
				post = postCtx.(*src.Post)
			}
			reqData.PostData = post
		case "posts":
			posts, _ := DbGetAllPosts()
			reqData.PostsData = posts
		case "replyTo":
			replyTo := src.GetParam(r.URL.String(), "reply_to")
			reply := &src.Reply{ReplyTo: replyTo}
			reqData.ReplyData = reply
		}
	}
	client := src.GetClientData(r)
	reqData.ClientData = client

	return reqData
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	src.Log(r)

	reqData := GetRequestDataC([]string{"posts"}, r)

	t, _ := template.ParseFiles("templates/base.html", "templates/index.html")
	err := t.Execute(w, reqData)
	src.Catch(err)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	reqData := GetRequestDataC([]string{"post"}, r)

	t, _ := template.ParseFiles("templates/base.html", "templates/post.html")
	err := t.Execute(w, reqData)
	src.Catch(err)
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	reqData := GetRequestDataC([]string{"replyTo"}, r)

	t, _ := template.ParseFiles("templates/base.html", "templates/create.html")
	err := t.Execute(w, reqData)
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
