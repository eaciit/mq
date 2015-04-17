package msg

import (
	"time"
)

type MqMsg struct {
	MsgType    string
	Value      interface{}
	Created    time.Time
	LastAccess time.Time
	Expiry     time.Duration
}
