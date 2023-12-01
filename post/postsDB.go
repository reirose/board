package post

import (
	"board/src"
	"fmt"
	"strconv"
)

func DbCheckForId(id int) bool {
	var idC int

	query, err := src.Database.Prepare("select id from posts where id = ?")
	if err != nil {
		return false
	}

	result := query.QueryRow(id)

	return result.Scan(&idC) == nil
}

func DbGetLastId() (int, error) {
	var data int
	query, err := src.Database.Prepare("select id from posts order by id desc limit 1")

	if err != nil {
		return 0, err
	}

	res := query.QueryRow()
	err = res.Scan(&data)

	if err != nil {
		fmt.Println(data, err)
		return 0, err
	}

	return data, nil
}

func DbGetAllPosts() ([]*src.Post, error) {
	query, err := src.Database.Prepare("select id, content, published_at, parent_id from posts limit 100")

	if err != nil {
		return nil, err
	}

	result, err := query.Query()
	if err != nil {
		return nil, err
	}

	posts := make([]*src.Post, 0)

	for result.Next() {
		data := new(src.Post)
		err := result.Scan(
			&data.ID,
			&data.Content,
			&data.PublishedAt,
			&data.ParentID,
		)
		if err != nil {
			return nil, err
		}

		childIDs := make([]int, 0)

		query, err = src.Database.Prepare("select id from posts where parent_id = ?")
		if err != nil {
			return nil, err
		}

		res, err := query.Query(data.ID)
		if err != nil {
			return nil, err
		}

		for res.Next() {
			dat := new(int)

			err = res.Scan(&dat)
			if err != nil {
				return nil, err
			}

			childIDs = append(childIDs, *dat)
		}

		data.ChildrenIDs = childIDs

		for _, v := range childIDs {
			replyPost := new(src.PostReply)
			postIDString := strconv.Itoa(v)
			post, err := DbGetPost(postIDString)
			src.Catch(err)
			replyPost.ID = post.ID
			replyPost.Content = post.Content
			replyPost.PublishedAt = post.PublishedAt
			replyPost.ChildrenIDs = post.ChildrenIDs

			data.Children = append(data.Children, replyPost)
		}

		if !DbCheckForId(data.ParentID) {
			data.ParentID = 0
		}

		posts = append(posts, data)
	}

	src.Catch(query.Close())
	return posts, nil
}

func DbGetPost(postID string) (*src.Post, error) {
	query, err := src.Database.Prepare("select id, content, published_at, parent_id from posts where id = ?")
	if err != nil {
		return nil, err
	}

	result := query.QueryRow(postID)

	data := new(src.Post)
	err = result.Scan(&data.ID, &data.Content, &data.PublishedAt, &data.ParentID)
	if err != nil {
		return nil, err
	}

	if !DbCheckForId(data.ParentID) {
		data.ParentID = 0
	}

	childs := make([]int, 0)

	query, err = src.Database.Prepare("select id from posts where parent_id = ?")
	if err != nil {
		return nil, err
	}

	res, err := query.Query(postID)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		data := new(int)
		err = res.Scan(&data)
		if err != nil {
			return nil, err
		}

		childs = append(childs, *data)
	}

	data.ChildrenIDs = childs

	src.Catch(query.Close())
	return data, nil
}

func DbCreatePost(post *src.Post) error {
	query, err := src.Database.Prepare("insert or replace into posts(content, published_at, parent_id) values (?,?,?)")
	if err != nil {
		return err
	}

	if _, err = query.Exec(post.Content, post.PublishedAt, post.ParentID); err != nil {
		return err
	}

	src.Catch(query.Close())
	return nil
}

func DbDeletePost(id int) error {
	query, err := src.Database.Prepare("delete from posts where id=?")
	if err != nil {
		return err
	}

	if _, err = query.Exec(id); err != nil {
		return err
	}

	src.Catch(query.Close())
	return nil
}
