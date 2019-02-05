package db

import _ "github.com/go-sql-driver/mysql"

func BuildMysqlDbGateway(dsnConnString string, maxConnAttempts int) (*DbGateway, error) {
	return BuildDbGateway("mysql", dsnConnString, maxConnAttempts)
}
