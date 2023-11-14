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

func getParam(w http.ResponseWriter, r_url string, req_param string) (string, error) {
	u, err := url.Parse(r_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	
	params := u.Query()
	res := params.Get(req_param)
	return res, err
}
