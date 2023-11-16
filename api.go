package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

func apiAnswer(r *http.Request) (*APIResponse, error) {
	response := new(APIResponse)
	userIds := make([]string, 0)
	userList := make([]*User, 0)
	postsList := make([]*Post, 0)

	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		catch(err)
		return nil, err
	}

	userIds = parsedURL.Query()["user_id"]

	for _, id := range userIds {
		user, err := dbGetUser(id)
		if err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}

	postsIds := parsedURL.Query()["post_id"]
	if postsIds[0] == "latest" {
		endIndex, err := dbGetLastId() // TODO: make sth like dbGetFirstAndLast()
		if err != nil {				   // i am 100% sure it can be done better but im too lazy lol
			return nil, err
		}

		length, err := dbGetAllPosts()
		if err != nil {
			return nil, err
		}
		if len(length) <= 50 {
			length = make([]*Post, 0)
		}

		for id := len(length); id <= endIndex; id++ {
			post, err := dbGetPost(strconv.Itoa(id))
			if err != nil {
				continue
			}
			postsList = append(postsList, post)
		}
	} else {
		for _, id := range postsIds {
			post, err := dbGetPost(id)
			if err != nil {
				continue
			}
			postsList = append(postsList, post)
		}
	}

	response.Version = "v1"
	response.Users = userList
	response.Posts = postsList
	jsonResp, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	response.JSON = jsonResp

	return response, nil
}

func APIAnswer(w http.ResponseWriter, r *http.Request) {
	log(r)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	data, err := apiAnswer(r)
	catch(err)

	n, err := w.Write(data.JSON)
	w.WriteHeader(n)
	catch(err)
}