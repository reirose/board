package user

import (
	"board/src"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
)

func DbGetUserByCookie(cs []*http.Cookie) (*src.User, error) {
	q, err := src.Database.Prepare("select user_id from users where token = ?")
	src.Catch(err)
	var c *http.Cookie

	for _, cc := range cs {
		if cc.Name == "user-token" {
			c = cc
		}
	}

	if c == nil {
		return nil, errors.New("cookie not found")
	}

	var userId string
	res := q.QueryRow(c.Value)
	err = res.Scan(&userId)

	resUser, err := DbGetUser(userId)
	if err != nil {
		return nil, err
	}

	return resUser, nil
}

func DbEncodeString(s string) *string {
	encSha256 := sha256.New()

	encSha256.Write([]byte(s))

	encString := fmt.Sprintf("%x", encSha256.Sum(nil))

	return &encString
}

func DbCheckEq(inputString string, dbPassword string) bool {
	encodedInputString := sha256.New()
	encodedInputString.Write([]byte(inputString))

	inputString = fmt.Sprintf("%x", encodedInputString.Sum(nil))

	return inputString == dbPassword
}

func DbGetUserById(uid int) (*src.User, error) {
	q, err := src.Database.Prepare("select user_id from users where id = ?")
	if err != nil {
		return nil, err
	}

	res := q.QueryRow(uid)
	data := new(string)
	src.Catch(res.Scan(&data))
	src.Catch(q.Close())
	return DbGetUser(*data)
}

func DbGetUser(user_id string) (*src.User, error) {
	q, err := src.Database.Prepare("select id, user_id, password, role, token from users where user_id = ?")
	if err != nil {
		return nil, err
	}

	res := q.QueryRow(user_id)

	data := new(src.User)

	err = res.Scan(
		&data.ID,
		&data.UserID,
		&data.Password,
		&data.Role,
		&data.Token,
	)
	if err != nil {
		return nil, err
	}

	src.Catch(q.Close())
	return data, nil
}

func DbRegisterUser(user *src.User) error {
	q, err := src.Database.Prepare("insert into users(user_id, password, role, token) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	if _, err = q.Exec(user.UserID, user.Password, user.Role, user.Token); err != nil {
		return err
	}

	src.Catch(q.Close())
	return nil
}

func DbRemoveUser(id int) error {
	q, err := src.Database.Prepare("delete from users where id = ?")
	if err != nil {
		return err
	}

	if _, err = q.Exec(id); err != nil {
		return nil
	}

	src.Catch(q.Close())
	return nil
}

func DbRevokeUser(user *src.User) error {
	q, err := src.Database.Prepare("update users set token = ? where id = ?")
	if err != nil {
		return err
	}

	if _, err = q.Exec(user.Token, user.ID); err != nil {
		return err
	}

	src.Catch(err)
	return nil
}
