package msg

import (
	"compress/gzip"
	"encoding/gob"
	"os"
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
	Size       int64
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

func BuildKey(owner string, table string, key string) string {
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

func (m *MqMsg) LoadFromFile(filename string) error {

	fi, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return err
	}
	defer fz.Close()

	decoder := gob.NewDecoder(fz)
	err = decoder.Decode(m)
	if err != nil {
		return err
	}

	return nil
}

func (m *MqMsg) SaveToFile(filename string) error {

	fi, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fi.Close()

	fz := gzip.NewWriter(fi)
	defer fz.Close()

	encoder := gob.NewEncoder(fz)
	err = encoder.Encode(m)
	if err != nil {
		return err
	}

	return nil
}
