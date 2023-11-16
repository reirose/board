package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	log(r)
	user := r.Context().Value("user").(*User)

	t, _ := template.ParseFiles("templates/base.html", "templates/user.html")
	err := t.Execute(w, user)
	catch(err)
}

func RegUser(w http.ResponseWriter, r *http.Request) {
	log(r)
	var err error
	preReg := new(PreRegUser)
	preReg.UserID = uuid.New().String()
	preReg.Role = getParam(r.URL.String(), "role")

	t, _ := template.ParseFiles("templates/base.html", "templates/register.html")
	err = t.Execute(w, preReg)
	catch(err)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	log(r)
	fmt.Println(r.Form)
	userId := r.FormValue("user_id")
	password, err := dbEncodeString(r.FormValue("password"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	role := r.FormValue("role")

	fmt.Println(userId, *password, role)

	user := &User{
		UserID:   userId,
		Password: *password,
		Role:     role,
	}

	catch(dbRegisterUser(user))
	http.Redirect(w, r, "/", http.StatusFound)
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	log(r)
	user := r.Context().Value("user").(*User)
	err := dbRemoveUser(user.ID)
	catch(err)

	http.Redirect(w, r, "/", http.StatusFound)
}
