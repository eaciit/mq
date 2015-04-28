package client

import (
	. "github.com/eaciit/mq/helper"
	. "github.com/eaciit/mq/msg"
	"net/rpc"
	"time"
)

type MqClient struct {
	connection *rpc.Client
	ClientInfo *ClientInfo
}

type ClientInfo struct {
	Username   string
	Password   string
	Role       string
	IsLoggedIn bool
	LastLogin  time.Time
}

func NewMqClient(dsn string, timeout time.Duration) (*MqClient, error) {
	rpcClient, err := rpc.Dial("tcp", dsn)
	if err != nil {
		return nil, err
	}
	ci := ClientInfo{}
	ci.IsLoggedIn = false

	return &MqClient{rpcClient, &ci}, nil
}

// func (c *MqClient) SetClientInfo() (bool, error) {
// 	c.clientInfo = ci
// 	return false, nil
// }

func (c *MqClient) Close() {
	c.connection.Close()
}

func (c *MqClient) Call(op string, key interface{}) (*MqMsg, error) {
	result := MqMsg{}
	err := c.connection.Call("MqRPC."+op, key, &result)
	return &result, err
}

func (c *MqClient) CallInc(op string, data string, key string) (*MqMsg, error) {
	result := MqMsg{} //
	k := map[string]interface{}{
		"data": data,
		"key":  key,
	}
	err := c.connection.Call("MqRPC.Inc", k, &result)
	return &result, err
}

func (c *MqClient) CallToLogin(key MqMsg) (*MqMsg, error) {
	result := MqMsg{}
	ci := ClientInfo{}
	err := c.connection.Call("MqRPC.ClientLogin", key, &result)
	if result.Value != "0" {
		//login success
		ci.IsLoggedIn = true
		ci.Username = key.Key
		ci.Password = key.Value.(string)
		ci.Role = result.Value.(string)
		ci.LastLogin = time.Now()
	} else {
		ci.IsLoggedIn = false
		ci.Username = ""
		ci.Password = ""
		ci.Role = ""
	}
	result.Value = ci
	c.ClientInfo = &ci
	return &result, err
}

func (c *MqClient) CallDecode(op string, key interface{}, resultPointer interface{}) error {
	result := MqMsg{}
	err := c.connection.Call("MqRPC."+op, key, &result)
	if err != nil {
		return err
	}
	Decode(result.Value.([]byte), resultPointer)
	return nil
}

func (c *MqClient) CallString(op string, key interface{}) (string, error) {
	result := MqMsg{}
	err := c.connection.Call("MqRPC."+op, key, &result)
	if err != nil {
		return "", err
	}
	return result.Value.(string), nil
}

func (c *MqClient) CallToLog(op string, key interface{}) (*MqMsg, error) {
	result := MqMsg{}
	err := c.connection.Call("MqRPC."+op, key, &result)
	return &result, err
}
