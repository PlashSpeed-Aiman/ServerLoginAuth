package model

import (
	"ServerLoginAuth/model/DB"
	"ServerLoginAuth/utils/crypto_utils"
	"context"
	"fmt"
	"log"
	"os"
)

type User struct {
	Username int    `json:"username",binding:"required"`
	Password string `json:"password",binding:"required"`
}

func CheckUserExist(username int) bool {
	var user User
	r_chan := make(chan error)
	go func() {
		err := DB.DBConn.QueryRow(context.Background(), "SELECT * FROM user_details WHERE USERNAME=$1", username).Scan(&user.Username, &user.Password)
		r_chan <- err
	}()
	if err := <-r_chan; err != nil {
		fmt.Fprintf(os.Stderr, "USER DOES NOT EXIST: %v\n", err)
		return false
	}
	log.Printf("USER EXISTS")
	return true
}
func CheckUser(username int, password string) bool {
	// DO NOT use this salt value; generate your own random salt. 8 bytes is
	// a good length.

	//DATABASE QUERY
	var user User
	saltedPass := crypto_utils.Salt_the_earth(password)
	//defer conn.Close()
	err_chan := make(chan error, 5)
	go func() {
		err := DB.DBConn.QueryRow(context.Background(), "SELECT * FROM user_details WHERE USERNAME=$1", username).Scan(&user.Username, &user.Password)
		err_chan <- err
	}()

	if err := <-err_chan; err != nil {
		fmt.Fprintf(os.Stderr, "USER DOES NOT EXIST: %v\n", err)
		return false
	}

	if user.Password != saltedPass {
		return false
	}

	return true
}

func RegisterUser(username string, password string) (bool, error) {
	saltedPass := crypto_utils.Salt_the_earth(password)
	r_chan := make(chan error, 5)
	//don't know if this improves performances
	go func() {
		_, err := DB.DBConn.Exec(context.Background(), "INSERT INTO user_details(username,password) values ($1,$2)", username, saltedPass)
		r_chan <- err
	}()
	if err := <-r_chan; err != nil {
		fmt.Fprintf(os.Stderr, "USER DOES NOT EXIST: %v\n", err)
		return false, err
	}

	return true, nil
}
