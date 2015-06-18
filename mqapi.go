package main

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"time"

	. "github.com/eaciit/mq/client"
	. "github.com/eaciit/mq/msg"
	"github.com/go-martini/martini"
)

type TokenData struct {
	Data struct {
		Token string    `json:"token"`
		Valid time.Time `json:"valid"`
	} `json:"data"`
	Result  string `json:"result"`
	Message string `json:"message"`
}

const (
	tokenLength = 32
)

func main() {
	m := martini.Classic()
	m.Get("/api/gettoken/username=(?P<name>[a-zA-Z0-9]+)&password=(?P<password>[a-zA-Z0-9]+)", GetToken)
	m.RunOnAddr(":8090")
}

func GetToken(params martini.Params) string {
	var result string
	auth, e := Auth(params["name"], params["password"])
	if auth && e == nil {
		result, _ = ResultValue("", true)
	} else {
		result, _ = ResultValue("Not authenticated user", false)
	}
	return result
}

func Auth(username, password string) (bool, error) {
	c, _ := NewMqClient("127.0.0.1:7890", time.Second*10)
	isLoggedIn := false
	// ActiveUser := ClientInfo{}
	// Role := ""
	msg := MqMsg{Key: username, Value: password}
	i, e := c.CallToLogin(msg)
	if e != nil {
		return false, e
	}
	// fmt.Println(i)
	if i.Value.(ClientInfo).IsLoggedIn {
		isLoggedIn = true
		// Role = i.Value.(ClientInfo).Role
		// ActiveUser = i.Value.(ClientInfo)
	}
	return isLoggedIn, nil
}
func ResultValue(message string, isAuthSuccess bool) (string, error) {
	token, _ := GenerateRandomString(tokenLength)
	now := time.Now()
	then := now.Add(10 * time.Minute)
	result := new(TokenData)
	if isAuthSuccess {
		result.Result = "OK"
		result.Message = ""
		result.Data.Token = token
		result.Data.Valid = then
	} else {
		result.Result = "error"
		result.Message = message
		result.Data.Token = ""
		result.Data.Valid = now
	}

	b, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(randInt(0, 100))
	}
	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
