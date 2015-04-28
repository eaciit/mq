package msg

import (
	"strings"
	"time"
)

type MqMsg struct {
	Key        string
	Value      interface{}
	Created    time.Time
	LastAccess time.Time
	Expiry     time.Duration
	Owner      string // user
	Duration   int64
	Table      string
	Permission string
}

func (msg *MqMsg) SetDefaults(m *MqMsg) {

	msg.Table = ""
	msg.Owner = "public"
	msg.Duration = 0
	msg.Permission = "666"

	if strings.TrimSpace(m.Table) != "" {
		msg.Table = m.Table
	}
	if strings.TrimSpace(m.Owner) != "" {
		msg.Owner = m.Owner
	}
	if m.Duration > 0 {
		msg.Duration = m.Duration
	}
	if strings.TrimSpace(m.Permission) != "" {
		msg.Permission = m.Permission
	}
}

func (msg *MqMsg) BuildKey(owner string, table string, key string) string {
	genKey := ""
	if strings.TrimSpace(owner) != "" {
		genKey = strings.TrimSpace(owner) + "|"
	} else {
		genKey = "public|"
	}
	if table != "" {
		genKey = genKey + strings.TrimSpace(table) + "|"

	}
	if key != "" {
		genKey = genKey + strings.TrimSpace(key)
	}

	return genKey
}
