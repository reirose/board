package src

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func Log(r *http.Request) {
	t := time.Now()
	fmt.Printf("%d:%d:%d - %s:%d\n", t.Hour(), t.Minute(), t.Second(), r.Method, r.ContentLength)
}

func GetParam(rUrl string, reqParam string) string {
	u, err := url.Parse(rUrl)
	if err != nil {
		return ""
	}

	res := u.Query().Get(reqParam)
	for _, i := range u.Query() {
		fmt.Println(i)
	}
	return res
}

func Catch(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

//func CheckCookie(r *http.Request, token string) (bool, error) {
//	c, err := r.Cookie("token")
//	if err != nil {
//		return false, err
//	}
//
//	if c.Value != token {
//		return false, nil
//	}
//
//	return true, nil
//}
