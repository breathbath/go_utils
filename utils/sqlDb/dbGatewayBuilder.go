package db

func BuildDbGateway(driverName, dsnConnString string, maxConnAttempts int) (*DbGateway, error) {
	conn, err := NewDb(dsnConnString, driverName)
	if err != nil {
		return nil, err
	}

	err = ValidateConn(conn, maxConnAttempts)
	if err != nil {
		return nil, err
	}

	return NewDbGateway(conn), nil
}
