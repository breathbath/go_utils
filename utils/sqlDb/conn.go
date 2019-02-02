package db

import (
	"github.com/breathbath/go_utils/utils/connections"
	"github.com/breathbath/go_utils/utils/io"
	"github.com/jmoiron/sqlx"
)

//LoggingFunc global variable func to add custom logging output
var LoggingFunc = func(msg string, err error) {
	if err != nil {
		io.OutputError(err, "", msg)
	} else {
		io.OutputInfo("", msg)
	}
}

//NewDb creates db connection with multiple retries in case if mysql is not immediately available
func NewDb(dsnConnString, sqlDriverName string, maxConnAttempts int) (*sqlx.DB, error) {
	return sqlx.Open(sqlDriverName, dsnConnString)
}

//ValidateConn pings the opened connection and fails if conn details are wrong
func ValidateConn(conn *sqlx.DB, maxConnAttempts int) error {
	_, err := connections.WaitForConnection(
		maxConnAttempts,
		"sql_db",
		func() (interface{}, error) {
			return nil, conn.Ping()
		},
		LoggingFunc,
	)

	return err
}
