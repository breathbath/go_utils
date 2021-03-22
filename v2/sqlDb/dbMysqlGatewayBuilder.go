package db

func BuildMysqlDBGateway(dsnConnString string, maxConnAttempts int) (*Gateway, error) {
	return BuildDBGateway("mysql", dsnConnString, maxConnAttempts)
}
