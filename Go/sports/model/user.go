package model

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func UserReg(username, password string) (int64, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	encodePwd := string(hash)
	data := User{
		Username:    username,
		Password:    encodePwd,
		State:       1,
		LastLoginIp: "127.0.0.1",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	query, _, _ := dialect.Insert("user").Rows(&data).ToSQL()
	ret, err := meta.MerchantDB.Exec(query)
	if err != nil {
		fmt.Println("insert user failed query", query)
		fmt.Println("insert user failed error", err.Error())
		return 0, err
	}

	id, _ := ret.LastInsertId()

	return id, nil
}

func UserExist(username string) int {

	var id int
	sql, _, _ := dialect.From("user").Select(goqu.COUNT("*").As("id")).Where(goqu.Ex{
		"username": username,
	}).ToSQL()
	//sql := fmt.Sprintf("SELECT COUNT(*) As id FROM user WHERE (username = '%s'", username)
	meta.MerchantDB.Get(&id, sql)
	return id
}

func UserLogin(username, password string) (UserLoginResp, error) {

	rsp := UserLoginResp{}

	return rsp, nil
}
