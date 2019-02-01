package db

import (
	"github.com/breathbath/go_utils/utils/connections"
	"github.com/breathbath/go_utils/utils/io"
	"github.com/jmoiron/sqlx"
)

func NewDb(dsnConnString string, maxConnAttempts int) (*sqlx.DB, error) {
	resource, err := connections.WaitForConnection(
		int(maxConnAttempts),
		"MYSQL",
		func() (interface{}, error) {
			return sqlx.Open("mysql", dsnConnString)
		},
		func(msg string, err error) {
			if err != nil {
				io.OutputError(err, "", msg)
			} else {
				io.OutputInfo("", msg)
			}
		},
	)
	if err != nil {
		return nil, err
	}

	return resource.(*sqlx.DB), nil
}
