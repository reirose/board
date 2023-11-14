package main

import (
	"crypto/md5"
	"fmt"
	"io"
)

func dbEncodeString(s string) (*string, error) {
	enc_md5 := md5.New()
	_, err := io.WriteString(enc_md5, s)
	if err != nil {
		return nil, err
	}

	enc_string := fmt.Sprintf("%x", enc_md5.Sum(nil))

	return &enc_string, nil
}

// func dbCheckEq(enc_s1 string, s2 string) (bool, error) {
// 	enc_s2 := md5.New()
// 	_, err := io.WriteString(enc_s2, s2)
// 	if err != nil {
// 		return false, err
// 	}

// 	encoded_s2 := fmt.Sprintf("%x", enc_s2.Sum(nil))

// 	return encoded_s2 == enc_s1, nil
// }

// func dbGetLastID() (string, error) {
// 	query, err := db.Prepare("select user_id from users order by id desc limit 1")
// 	if err != nil {
// 		return "", err
// 	}

// 	var data string
// 	result := query.QueryRow()
// 	err = result.Scan(data)
// 	if err != nil {
// 		return "", nil
// 	}
	
// 	defer query.Close()
// 	return data, nil
// }

func dbGetUser(UserID string) (*User, error) {
	q, err := db.Prepare("select id, user_id, password, role from users where id = ?")
	if err != nil {
		return nil, err
	}

	res := q.QueryRow(UserID)
	if err != nil {
		return nil, err
	}

	data := new(User)
	
	err = res.Scan(
		&data.ID,
		&data.UserID,
		&data.Password,
		&data.Role,
	)
	if err != nil {
		return nil, err
	}

	defer q.Close()
	return data, nil
}

func dbRegisterUser(user *User) error {
	q, err := db.Prepare("insert into users(user_id, password, role) values (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = q.Exec(user.UserID, user.Password, user.Role)
	if err != nil {
		return err
	}

	defer q.Close()
	return nil
}

func dbRemoveUser(id int) error {
	q, err := db.Prepare("delete from users where id = ?")
	if err != nil {
		return nil
	}

	_, err = q.Exec(id)
	if err != nil {
		return nil
	}
	
	defer q.Close()
	return nil
}
