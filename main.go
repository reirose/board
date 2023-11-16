package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var router *chi.Mux
var db *sql.DB

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "postID")
		post, err := dbGetPost(postID)
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
		userID := chi.URLParam(r, "userID")
		user, err := dbGetUser(userID)
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
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)

	var err error
	db, err = connect()
	catch(err)

	router.Use(ChangeMethod)

	// index methods
	router.Get("/", GetAllPosts)

	// /create/ methods
	router.Route("/create", func(r chi.Router) {
		r.Get("/", NewPost)
		r.Post("/", CreatePost)
	})

	// /post/ methods
	router.Route("/post", func(r chi.Router) {
		r.Get("/", ErrorNotFound)
		r.Route("/{postID}", func(r chi.Router) {
			r.Use(PostCtx)
			r.Get("/", GetPost)
			r.Delete("/", DeletePost)
			r.Get("/delete", DeletePost)
		})
	})

	// /user/ methods
	router.Route("/user", func(r chi.Router) {
		r.Get("/", ErrorNotFound)
		r.Route("/{userID}", func(r chi.Router) {
			r.Use(UserCtx)
			r.Get("/", GetUser)
			r.Get("/remove", RemoveUser)
		})
		r.Route("/reg", func(r chi.Router) {
			r.Get("/", RegUser)
			r.Post("/", AddUser)
		})
	})

	// API
	router.Route("/api", func(r chi.Router) {
		r.Get("/", APIAnswer)
	})

	fileServer := http.FileServer(http.Dir("./assets/"))
	router.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	fmt.Println("Listening and serving @ localhost:3000")
	catch(http.ListenAndServe(":3000", router))
}
