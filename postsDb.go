package main

import "fmt"

func dbCheckForId(id int) bool {
	var idC int

	query, err := db.Prepare("select id from posts where id = ?")
	if err != nil {
		return false
	}

	result := query.QueryRow(id)

	return result.Scan(&idC) == nil
}

func dbGetLastId() (int, error) {
	var data int
	query, err := db.Prepare("select id from posts order by id desc limit 1")

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

func dbGetAllPosts() ([]*Post, error) {
	query, err := db.Prepare("select id, content, published_at, parent_id from posts limit 100")

	if err != nil {
		return nil, err
	}

	result, err := query.Query()
	if err != nil {
		return nil, err
	}

	posts := make([]*Post, 0)

	for result.Next() {
		data := new(Post)
		err := result.Scan(
			&data.ID,
			&data.Content,
			&data.PublishedAt,
			&data.ParentID,
		)
		if err != nil {
			return nil, err
		}
		
		childs := make([]int, 0)
		
		query, err = db.Prepare("select id from posts where parent_id = ?")
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
		
			childs = append(childs, *dat)
		}
		
		data.ChildrenIDs = childs

		if !dbCheckForId(data.ParentID) {
			data.ParentID = 0
		}
		
		posts = append(posts, data)
	}


	catch(query.Close())
	return posts, nil
}

func dbGetPost(postID string) (*Post, error) {
	query, err := db.Prepare("select id, content, published_at, parent_id from posts where id = ?")
	if err != nil {
		return nil, err
	}

	result := query.QueryRow(postID)

	data := new(Post)
	err = result.Scan(&data.ID, &data.Content, &data.PublishedAt, &data.ParentID)
	if err != nil {
		return nil, err
	}

	if !dbCheckForId(data.ParentID) {
		data.ParentID = 0
	}

	childs := make([]int, 0)

	query, err = db.Prepare("select id from posts where parent_id = ?")
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

	catch(query.Close())
	return data, nil
}

func dbCreatePost(post *Post) error {
	query, err := db.Prepare("insert or replace into posts(content, published_at, parent_id) values (?,?,?)")
	if err != nil {
		return err
	}

	_, err = query.Exec(post.Content, post.PublishedAt, post.ParentID)
	if err != nil {
		return err
	}

	catch(query.Close())
	return nil
}

func dbDeletePost(id int) error {
	query, err := db.Prepare("delete from posts where id=?")
	if err != nil {
		return err
	}

	_, err = query.Exec(id)
	if err != nil {
		return err
	}

	catch(query.Close())
	return nil
}
