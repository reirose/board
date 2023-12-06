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

func GetClientData(r *http.Request) *ClientData {
	cookies := r.Cookies()
	client := new(ClientData)
	for _, v := range cookies {
		if v.Name == "token" {
			client.UserToken = v.Value
		}
	}
	return client
}

func GetParam(rUrl string, reqParam string) string {
	u, err := url.Parse(rUrl)
	if err != nil {
		return ""
	}

	res := u.Query().Get(reqParam)
	return res
}

func Catch(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
