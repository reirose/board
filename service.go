package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func log(r *http.Request) {
	t := time.Now()
	fmt.Printf("%d:%d:%d - %s:%d\n", t.Hour(), t.Minute(), t.Second(), r.Method, r.ContentLength)
}

func getParam(r_url string, req_param string) (string) {
	u, err := url.Parse(r_url)
	if err != nil {
		return ""
	}

	res := u.Query().Get(req_param)
	for _, i := range u.Query() {
		fmt.Println(i)
	}
	return res
}
