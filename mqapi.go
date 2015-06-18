package main

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"time"

	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/helper"
	. "github.com/eaciit/mq/msg"
	"github.com/go-martini/martini"
)

type User struct {
	Username string
	Token    string
	Valid    time.Time
}

type TokenData struct {
	Token string    `json:"token"`
	Valid time.Time `json:"valid"`
}

const (
	tokenLength = 32
	expiredTime = 15
)

var (
	users []User
	user  User
)

func main() {
	m := martini.Classic()
	m.Post("/api/gettoken/username=(?P<username>[a-zA-Z0-9]+)&password=(?P<password>[a-zA-Z0-9]+)", GetToken)
	m.Post("/api/checktoken/tokenkey=(?P<tokenkey>[a-zA-Z0-9=]+)", CheckToken)
	m.RunOnAddr(":8090")
}

func CheckToken(w http.ResponseWriter, params martini.Params) {
	data := TokenData{}
	isTokenExist := false
	for _, v := range users {
		if params["tokenkey"] == v.Token {
			isTokenExist = true
			data.Token = v.Token
			data.Valid = v.Valid
			PrintJSON(w, true, data, "")
			break
		}
	}
	if !isTokenExist {
		PrintJSON(w, false, "", "token doesn't exist")
	}

}

func GetToken(w http.ResponseWriter, params martini.Params) {
	username := params["username"]
	password := params["password"]
	auth, e := Auth(username, password)
	if auth && e == nil {
		token := GenerateRandomString(tokenLength)
		valid := time.Now().Add(expiredTime * time.Minute)
		user.Username = username
		user.Token = token
		user.Valid = valid
		if CheckExistedUser(username, valid) {
			if CheckExpiredToken(user) {
				user.Token = token
				user.Valid = valid
				UpdateTokenAndValidTime(user)
			}
		} else {
			users = append(users, user)
		}
		data := TokenData{}
		data.Token = user.Token
		data.Valid = user.Valid
		PrintJSON(w, true, data, "")
	} else {
		PrintJSON(w, false, "", "wrong username and password combination")
	}
}

func UpdateTokenAndValidTime(usr User) {
	for k, v := range users {
		if usr.Username == v.Username {
			users[k].Token = usr.Token
			users[k].Valid = usr.Valid
			break
		}
	}
}

func CheckExistedUser(username string, valid time.Time) bool {
	isExist := false
	for k, v := range users {
		if v.Username == username {
			users[k].Valid = valid
			user = users[k]
			isExist = true
			break
		}
	}
	return isExist
}

func CheckExpiredToken(usr User) bool {
	return usr.Valid.Before(time.Now())
}

func Auth(username, password string) (bool, error) {
	c, _ := NewMqClient("127.0.0.1:7890", time.Second*10)
	isLoggedIn := false
	msg := MqMsg{Key: username, Value: password}
	i, e := c.CallToLogin(msg)
	if e != nil {
		return false, e
	}
	if i.Value.(ClientInfo).IsLoggedIn {
		isLoggedIn = true
	}
	return isLoggedIn, nil
}

func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(randInt(0, 100))
	}
	return b
}

func GenerateRandomString(s int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
