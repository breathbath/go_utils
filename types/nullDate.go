package types

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
)

type NullTime struct {
	mysql.NullTime
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nt.Time)
}
