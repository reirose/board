package user

import (
	"crypto/sha256"
	"fmt"
	"github.com/reirose/board/src"
)

func DbEncodeString(s string) (*string, error) {
	encSha256 := sha256.New()

	encSha256.Write([]byte(s))

	encString := fmt.Sprintf("%x", encSha256.Sum(nil))

	return &encString, nil
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

func DbGetUser(UserID string) (*src.User, error) {
	q, err := src.Database.Prepare("select id, user_id, password, role, token from users where id = ?")
	if err != nil {
		return nil, err
	}

	res := q.QueryRow(UserID)
	if err != nil {
		return nil, err
	}

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
