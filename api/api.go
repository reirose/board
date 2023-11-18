package api

import (
	"encoding/json"
	postlib "github.com/reirose/board/post"
	"github.com/reirose/board/src"
	userlib "github.com/reirose/board/user"
	"net/http"
	"net/url"
	"strconv"
)

func apiGetInfo(r *http.Request) (*src.APIResponse, error) {
	response := new(src.APIResponse)
	userIds := make([]string, 0)
	userList := make([]*src.User, 0)
	postsList := make([]*src.Post, 0)

	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		src.Catch(err)
		return nil, err
	}

	userIds = parsedURL.Query()["user_id"]

	for _, id := range userIds {
		user, err := userlib.DbGetUser(id)
		if err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}

	postsIds := parsedURL.Query()["post_id"]
	if len(postsIds) != 0 {
		if postsIds[0] == "latest" {
			endIndex, err := postlib.DbGetLastId() // TODO: make sth like dbGetFirstAndLast()
			if err != nil {                        // I am 100% sure it can be done better but im too lazy lol
				return nil, err
			}

			length, err := postlib.DbGetAllPosts()
			if err != nil {
				return nil, err
			}

			if len(length) >= 50 {
				length = make([]*src.Post, 0)
			}

			for id := len(length); id <= endIndex; id++ {
				post, err := postlib.DbGetPost(strconv.Itoa(id))
				if err != nil {
					continue
				}
				postsList = append(postsList, post)
			}
		} else {
			for _, id := range postsIds {
				post, err := postlib.DbGetPost(id)
				if err != nil {
					continue
				}
				postsList = append(postsList, post)
			}
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

func GetInfo(w http.ResponseWriter, r *http.Request) {
	src.Log(r)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	data, err := apiGetInfo(r)
	src.Catch(err)

	n, err := w.Write(data.JSON)
	w.WriteHeader(n)
	src.Catch(err)
}