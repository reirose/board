package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func log(r *http.Request) {
	t := time.Now()
	ip, err := getIP(r)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("%d:%d:%d - %s:%d @ %s\n", t.Hour(), t.Minute(), t.Second(), r.Method, r.ContentLength, ip)
}

func getParam(w http.ResponseWriter, r_url string, req_param string) (string, error) {
	u, err := url.Parse(r_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	
	res := u.Query().Get(req_param)
	return res, nil
}

func getIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}
