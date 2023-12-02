package user

import (
	"fmt"
	"html/template"
	"net/http"

	"board/post"
	"board/src"

	"github.com/google/uuid"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	src.Log(r)

	//user := r.Context().Value("user").(*src.User)
	reqData := post.GetRequestDataC([]string{"user"}, r)

	t, _ := template.ParseFiles("templates/base.html", "templates/user.html")
	err := t.Execute(w, reqData)
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

	dbUser, err := DbGetUserByCookie(r.Cookies())
	if dbUser != nil {
		post.GetAllPosts(w, r)
		return
	}

	t, _ := template.ParseFiles("templates/base.html", "templates/login.html")
	err = t.Execute(w, "")
	src.Catch(err)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	userId := r.FormValue("user_id")

	password := r.FormValue("password")

	user, err := DbGetUser(userId)

	if user != nil && DbCheckEq(password, user.Password) {
		t, _ := template.ParseFiles("templates/base.html", "templates/login.html")
		header := fmt.Sprintf("token=%s; samesite=None; max-age=%d; secure=true", user.Token, 3600*24*15)
		fmt.Println(header)
		w.Header().Set("Set-Cookie", header)
		err := t.Execute(w, "false")
		src.Catch(err)
		return
	}

	t, _ := template.ParseFiles("templates/base.html", "templates/login.html")
	err = t.Execute(w, "true")
	src.Catch(err)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	userId := r.FormValue("user_id")
	password := DbEncodeString(r.FormValue("password"))
	unhashedToken := src.GenerateToken(userId)
	token := DbEncodeString(unhashedToken)

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
	user.Token = *DbEncodeString(user.Token)
	fmt.Println(user.Token)

	err := DbRevokeUser(user)
	src.Catch(err)

	http.Redirect(w, r, "/", http.StatusOK)
}
