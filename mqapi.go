package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
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
	Token string
	Valid time.Time
}

type PutData struct {
	Node        int
	Owner       string
	Valid, Size int64
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
	port := flag.Int("port", 8090, "Port of RCP call. Default is 1234")
	serverHost := flag.String("master", "127.0.0.1:7890", "Default master host")
	flag.Parse()

	client, _ := NewMqClient(*serverHost, time.Second*10)

	m := martini.Classic()
	m.Post("/api/gettoken/username=(?P<username>[a-zA-Z0-9]+)&password=(?P<password>[a-zA-Z0-9]+)", func(w http.ResponseWriter, params martini.Params) {
		GetToken(w, params, client)
	})
	m.Post("/api/checktoken/token=(?P<token>[a-zA-Z0-9=_-]+)", CheckToken)
	m.Get("/api/get/token=(?P<token>[a-zA-Z0-9]+)&key=(?P<key>[a-zA-Z0-9]+)", func(w http.ResponseWriter, params martini.Params) {
		Get(w, params, client)
	})
	m.Post("/api/put/token=(?P<token>[a-zA-Z0-9]+)&key=(?P<key>[a-zA-Z0-9]+)", func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		Put(w, r, params, client)
	})

	m.RunOnAddr(fmt.Sprint(":", *port))
}

func CheckToken(w http.ResponseWriter, params martini.Params) {
	data := TokenData{}
	isTokenExist := false
	for _, v := range users {
		if params["token"] == v.Token && !isTimeExpired(v.Valid) {
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

func GetToken(w http.ResponseWriter, params martini.Params, c *MqClient) {
	username := params["username"]
	password := params["password"]
	auth, e := Auth(username, password, c)
	if auth && e == nil {
		token := GenerateRandomString(tokenLength)
		valid := time.Now().Add(expiredTime * time.Minute)
		user.Username = username
		user.Token = token
		user.Valid = valid
		if CheckExistedUser(username, valid) {
			if isTimeExpired(user.Valid) {
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

func isTimeExpired(valid time.Time) bool {
	return valid.Before(time.Now())
}

func checkTokenValidity(token string) error {
	for _, u := range users {
		if u.Token == token {
			if isTimeExpired(u.Valid) {
				return errors.New("TOKEN NOT VALID: Token already expired request new token.")
			} else {
				return nil
			}
		}
	}

	return errors.New("TOKEN NOT VALID: Token not found.")
}

func Get(w http.ResponseWriter, params martini.Params, c *MqClient) {
	e := checkTokenValidity(params["token"])
	if e != nil {
		PrintJSON(w, false, "", e.Error())
	} else {
		result, err := c.Call("Get", "public|"+params["key"])
		if err != nil {
			PrintJSON(w, false, "", err.Error())
		} else {
			PrintJSON(w, true, result, "")
		}
	}
}

func Put(w http.ResponseWriter, r *http.Request, params martini.Params, c *MqClient) {
	e := checkTokenValidity(params["token"])
	if e != nil {
		PrintJSON(w, false, "", e.Error())
	} else {
		key := BuildKey("", "", params["key"])
		arg := MqMsg{Key: key, Value: r.FormValue("value")}

		item, err := c.Call("Set", arg)

		if err != nil {
			PrintJSON(w, false, "", err.Error())
		} else {
			result := PutData{Owner: item.Owner, Size: item.Size, Valid: item.Duration}
			PrintJSON(w, true, result, "")
		}
	}
}

func Auth(username, password string, c *MqClient) (bool, error) {
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
