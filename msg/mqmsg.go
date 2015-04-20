package msg

import (
	"time"
)

type MqMsg struct {
	Key        string
	Value      interface{}
	Created    time.Time
	LastAccess time.Time
	Expiry     time.Duration
}
