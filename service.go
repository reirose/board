package main

import (
	"fmt"
	"net/http"
	"time"
)

func log(r *http.Request) {
	t := time.Now()
	fmt.Printf("%d:%d:%d - %s:%d\n", t.Hour(), t.Minute(), t.Second(), r.Method, r.ContentLength)
}
