package db

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/breathbath/go_utils/v2/connections"
	testing2 "github.com/breathbath/go_utils/v2/testing"
	"github.com/stretchr/testify/assert"
)

func TestOpenConnection(t *testing.T) {
	returnErr := true
	fakeDB := FakeSQLDriver{ConnStr: "", ErrFuncProvider: func() error {
		if returnErr {
			returnErr = false
			return errors.New("first err")
		}
		return nil
	}}

	sql.Register("fake_driver_TestOpenConnection", &fakeDB)

	sqlX, err := NewDB("conn_str", "fake_driver_TestOpenConnection")
	assert.NoError(t, err)

	connections.WaitingConnectorIterativeDelayDuration = time.Microsecond
	output := testing2.CaptureOutput(func() {
		err = ValidateConn(sqlX, 2)
		assert.NoError(t, err)
	})

	assert.Equal(t, "conn_str", fakeDB.ConnStr)
	testing2.AssertLogText(t, "[ERROR] sql_db connection error: first err[INFO] Trying to reconnect to sql_db in 0 s", output)
}
