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

var users []User

func main() {
	port := flag.Int("port", 8090, "Port of RCP call. Default is 1234")
	serverHost := flag.String("master", "127.0.0.1:7890", "Default master host")
	flag.Parse()

	client, _ := NewMqClient(*serverHost, time.Second*10)

	m := martini.Classic()
	m.Post("/api/gettoken/username=(?P<username>[a-zA-Z0-9]+)&password=(?P<password>[a-zA-Z0-9]+)", func(w http.ResponseWriter, params martini.Params) {
		GetToken(w, params, client)
	})
	m.Get("/api/checktoken/token=(?P<token>[a-zA-Z0-9=_-]+)", CheckToken)
	m.Get("/api/get/token=(?P<token>[a-zA-Z0-9=_-]+)&key=(?P<key>[a-zA-Z0-9]+)", func(w http.ResponseWriter, params martini.Params) {
		Get(w, params, client)
	})
	m.Post("/api/put/token=(?P<token>[a-zA-Z0-9=_-]+)&key=(?P<key>[a-zA-Z0-9]+)", func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		Put(w, r, params, client)
	})

	m.RunOnAddr(fmt.Sprint(":", *port))
}

func CheckToken(w http.ResponseWriter, params martini.Params) {
	u, err := checkTokenValidity(params["token"])
	if err != nil {
		PrintJSON(w, false, "", "token doesn't exist")
	} else {
		PrintJSON(w, true, TokenData{Token: u.Token, Valid: u.Valid}, "")
	}
}

func GetToken(w http.ResponseWriter, params martini.Params, c *MqClient) {
	username := params["username"]
	password := params["password"]
	auth, e := Auth(username, password, c)
	if auth && e == nil {
		foundIndex := getUserIndexWithUsername(username)
		data := TokenData{}
		if foundIndex >= 0 {
			if isTimeExpired(users[foundIndex].Valid) {
				users[foundIndex].Token = GenerateRandomString(tokenLength)
				users[foundIndex].Valid = time.Now().Add(expiredTime * time.Minute)
			}
			data.Token = users[foundIndex].Token
			data.Valid = users[foundIndex].Valid
		} else {
			user := User{}
			user.Username = username
			user.Token = GenerateRandomString(tokenLength)
			user.Valid = time.Now().Add(expiredTime * time.Minute)
			users = append(users, user)

			data.Token = user.Token
			data.Valid = user.Valid
		}
		PrintJSON(w, true, data, "")
	} else {
		PrintJSON(w, false, "", "wrong username and password combination")
	}
}

func updateTokenAndValidTime(user *User) {
	user.Valid = time.Now().Add(expiredTime * time.Minute)
}

func getUserIndexWithUsername(username string) int {
	for i, u := range users {
		if u.Username == username {
			return i
		}
	}
	return -1 //-1 Means user not found
}

func isTimeExpired(valid time.Time) bool {
	return valid.Before(time.Now())
}

func checkTokenValidity(token string) (User, error) {
	for _, u := range users {
		if u.Token == token {
			if isTimeExpired(u.Valid) {
				return User{}, errors.New("TOKEN NOT VALID: Token already expired request new token.")
			} else {
				updateTokenAndValidTime(&u)
				return u, nil
			}
		}
	}

	return User{}, errors.New("TOKEN NOT VALID: Token not found.")
}

func Get(w http.ResponseWriter, params martini.Params, c *MqClient) {
	_, e := checkTokenValidity(params["token"])
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
	_, e := checkTokenValidity(params["token"])
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
	msg := MqMsg{Key: username, Value: password}
	i, e := c.CallToLogin(msg)
	if e != nil {
		return false, e
	}

	return i.Value.(ClientInfo).IsLoggedIn, nil
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
