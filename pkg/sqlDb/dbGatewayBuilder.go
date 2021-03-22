package db

func BuildDBGateway(driverName, dsnConnString string, maxConnAttempts int) (*Gateway, error) {
	conn, err := NewDB(dsnConnString, driverName)
	if err != nil {
		return nil, err
	}

	err = ValidateConn(conn, maxConnAttempts)
	if err != nil {
		return nil, err
	}

	return NewDBGateway(conn), nil
}
