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
	preReg.Role, err = getParam(w, r.URL.String(), "role")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println(preReg.UserID, preReg.Role)

	t, _ := template.ParseFiles("templates/base.html", "templates/register.html")
	err = t.Execute(w, preReg)
	catch(err)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	log(r)
	fmt.Println(r.Form)
	user_id := r.FormValue("user_id")
	password, err := dbEncodeString(r.FormValue("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	role := r.FormValue("role")

	fmt.Println(user_id, *password, role)

	user := &User{
		UserID: user_id,
		Password: *password,
		Role: role,
	}


	err = dbRegisterUser(user)
	catch(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	log(r)
	user := r.Context().Value("user").(*User)
	err := dbRemoveUser(user.ID)
	catch(err)

	http.Redirect(w, r, "/", http.StatusFound)
}
