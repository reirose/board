package main

import (
	"board/api"
	postlib "board/post"
	"board/src"
	userlib "board/user"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "postID")
		post, err := postlib.DbGetPost(postID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "post", post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
		user, err := userlib.DbGetUserById(userID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ChangeMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch method := r.PostFormValue("_method"); method {
			case http.MethodPut:
				fallthrough
			case http.MethodPatch:
				fallthrough
			case http.MethodDelete:
				r.Method = method
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	src.Router = chi.NewRouter()
	src.Router.Use(middleware.Recoverer)

	var err error
	src.Database, err = connect()
	src.Catch(err)

	src.Router.Use(ChangeMethod)

	// index methods
	src.Router.Get("/", postlib.GetAllPosts)

	// /create/ methods
	src.Router.Route("/create", func(r chi.Router) {
		r.Get("/", postlib.NewPost)
		r.Post("/", postlib.CreatePost)
	})

	// /post/ methods
	src.Router.Route("/post", func(r chi.Router) {
		r.Get("/", postlib.ErrorNotFound)
		r.Route("/{postID}", func(r chi.Router) {
			r.Use(PostCtx)
			r.Get("/", postlib.GetPost)
			r.Delete("/", postlib.DeletePost)
			r.Get("/delete", postlib.DeletePost)
		})
	})

	// /user/ methods
	src.Router.Route("/user", func(r chi.Router) {
		r.Get("/", postlib.ErrorNotFound)
		r.Route("/{userID}", func(r chi.Router) {
			r.Use(UserCtx)
			r.Get("/", userlib.GetUser)
			r.Get("/revoke", userlib.Revoke)
			r.Get("/remove", userlib.RemoveUser)
		})
	})

	src.Router.Route("/login", func(r chi.Router) {
		r.Get("/", userlib.LoginPage)
		r.Post("/", userlib.LoginUser)
	})

	src.Router.Route("/register", func(r chi.Router) {
		r.Get("/", userlib.RegUser)
		r.Post("/", userlib.AddUser)
	})

	// API
	src.Router.Route("/api", func(r chi.Router) {
		r.Get("/get-info", api.GetInfo)
	})

	fileServer := http.FileServer(http.Dir("./assets/"))
	src.Router.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	fmt.Println("Listening and serving @ localhost:3000")
	src.Catch(http.ListenAndServe(":3000", src.Router))
}
