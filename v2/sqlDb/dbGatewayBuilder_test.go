package db

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/breathbath/go_utils/v2/connections"
	testing2 "github.com/breathbath/go_utils/v2/testing"
	"github.com/stretchr/testify/assert"
)

func TestDbGatewayBuilder(t *testing.T) {
	mysqlFakeDriver := NewFakeSQLDriver()
	sql.Register("TestDbGatewayBuilder_driver", mysqlFakeDriver)

	_, err := BuildDBGateway("TestDbGatewayBuilder_driver", "some_conn_str", 1)
	assert.NoError(t, err)

	assert.Equal(t, "some_conn_str", mysqlFakeDriver.ConnStr)

	mysqlFakeDriver.ErrFuncProvider = func() error {
		return fmt.Errorf("connection error")
	}
	connections.WaitingConnectorIterativeDelayDuration = time.Microsecond

	output := testing2.CaptureOutput(func() {
		_, err := BuildDBGateway("TestDbGatewayBuilder_driver", "some_conn_str", 1)
		assert.EqualError(t, err, "Was not able to connect to sql_db")
	})

	testing2.AssertLogText(t, "[ERROR] sql_db connection error: connection error[INFO] Trying to reconnect to sql_db in 0 s", output)
}

func TestDbGatewayBuilderWithInvalidDriver(t *testing.T) {
	_, err := BuildDBGateway("someDriver", "some_conn_str", 1)
	assert.EqualError(t, err, "sql: unknown driver \"someDriver\" (forgotten import?)")
}
