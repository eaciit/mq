package client

import (
	"time"
	"net/rpc"
	. "github.com/eaciit/mq/msg"
)

type MqClient struct {
	connection *rpc.Client
}

func NewMqClient(dsn string, timeout time.Duration) (*MqClient, error) {
	rpcClient, err := rpc.Dial("tcp", dsn)
	if err != nil {
		return nil, err
	}
	return &MqClient{rpcClient}, nil
}

func (c *MqClient) Close() {
	c.connection.Close()
}

func (c *MqClient) Call(op string, key interface{}) (*MqMsg, error) {
	result := MqMsg{}
	err := c.connection.Call("MqRPC."+op, key, &result)
	return &result, err
}

func (c *MqClient) CallString(op string, key interface{}) (string, error) {
	result := MqMsg{}
	err := c.connection.Call("MqRPC."+op, key, &result)
	if err != nil {
		return "", err
	}
	return result.Value.(string), nil
}
