package db

import (
	"database/sql"
	"errors"
	"github.com/breathbath/go_utils/utils/connections"
	testing2 "github.com/breathbath/go_utils/utils/testing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOpenConnection(t *testing.T) {
	returnErr := true
	fakeDb := FakeSqlDriver{ConnStr: "", ErrFuncProvider: func() error {
		if returnErr {
			returnErr = false
			return errors.New("First err")
		}
		return nil
	}}

	sql.Register("fake_driver_TestOpenConnection", &fakeDb)

	sqlX, err := NewDb("conn_str", "fake_driver_TestOpenConnection", 1)
	assert.NoError(t, err)

	connections.WaitingConnectorIterativeDelayDuration = time.Microsecond
	output := testing2.CaptureOutput(func() {
		err = ValidateConn(sqlX, 2)
		assert.NoError(t, err)
	})

	assert.Equal(t, "conn_str", fakeDb.ConnStr)
	testing2.AssertLogText(t,"[ERROR] sql_db connection error: First err[INFO] Trying to reconnect to sql_db in 0 s", output)
}
