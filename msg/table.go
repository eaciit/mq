package msg

import (
	"time"
)

type MqTable struct {
	TableId    string
	Owner      string
	Created    time.Time
	LastAccess time.Time
	Expiry     time.Duration
	Items      map[string]interface{}
}

func NewTable(tableid string, owner string) *MqTable {
	ret := new(MqTable)
	ret.TableId = tableid
	ret.Owner = owner
	ret.Created = time.Now()
	return ret
}
