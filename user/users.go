package user

import (
	"board/src"
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	user := r.Context().Value("user").(*src.User)

	t, _ := template.ParseFiles("templates/base.html", "templates/user.html")
	err := t.Execute(w, user)
	src.Catch(err)
}

func RegUser(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	var err error
	preReg := new(src.PreRegUser)
	preReg.UserID = uuid.New().String()
	preReg.Role = src.GetParam(r.URL.String(), "role")

	t, _ := template.ParseFiles("templates/base.html", "templates/register.html")
	err = t.Execute(w, preReg)
	src.Catch(err)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	t, _ := template.ParseFiles("templates/base.html", "templates/login.html")
	fmt.Println(r.Cookies())
	err := t.Execute(w, "")
	src.Catch(err)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	userId := r.FormValue("user_id")

	password, err := DbEncodeString(r.FormValue("password"))
	src.Catch(err)

	user, err := DbGetUser(userId)

	if (DbCheckEq(*password, user.Password)) || (err != nil) {
		t, _ := template.ParseFiles("templates/base.html", "templates/login.html")
		err := t.Execute(w, "true")
		src.Catch(err)
		return
	}

	w.Header().Set("Set-Cookie", fmt.Sprintf("%s=%s; samesite=None; max-age=%d", "foo", "bar",
		3600*24*15))

	t, _ := template.ParseFiles("templates/base.html", "templates/login.html")
	err = t.Execute(w, "false")
	src.Catch(err)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	userId := r.FormValue("user_id")
	password, err := DbEncodeString(r.FormValue("password"))
	unhashedToken := src.GenerateToken(userId)
	token, err := DbEncodeString(unhashedToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	role := r.FormValue("role")

	user := &src.User{
		UserID:   userId,
		Password: *password,
		Role:     role,
		Token:    *token,
	}

	src.Catch(DbRegisterUser(user))
	http.Redirect(w, r, "/", http.StatusFound)
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	user := r.Context().Value("user").(*src.User)
	err := DbRemoveUser(user.ID)
	src.Catch(err)

	http.Redirect(w, r, "/", http.StatusFound)
}

func Revoke(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	user := r.Context().Value("user").(*src.User)
	user.Token = src.GenerateToken(user.UserID)

	err := DbRevokeUser(user)
	src.Catch(err)

	http.Redirect(w, r, "/", http.StatusOK)
}
