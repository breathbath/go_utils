package db

import (
	"database/sql"
	"fmt"
	"github.com/breathbath/go_utils/utils/connections"
	testing2 "github.com/breathbath/go_utils/utils/testing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDbGatewayBuilder(t *testing.T) {
	mysqlFakeDriver := NewFakeSqlDriver()
	sql.Register("TestDbGatewayBuilder_driver", mysqlFakeDriver)

	_, err := BuildDbGateway("TestDbGatewayBuilder_driver", "some_conn_str", 1)
	assert.NoError(t, err)

	assert.Equal(t, "some_conn_str", mysqlFakeDriver.ConnStr)

	mysqlFakeDriver.ErrFuncProvider = func() error {
		return fmt.Errorf("Connection error")
	}
	connections.WaitingConnectorIterativeDelayDuration = time.Microsecond

	output := testing2.CaptureOutput(func() {
		_, err := BuildDbGateway("TestDbGatewayBuilder_driver", "some_conn_str", 1)
		assert.EqualError(t, err, "Was not able to connect to sql_db")
	})

	testing2.AssertLogText(t,"[ERROR] sql_db connection error: Connection error[INFO] Trying to reconnect to sql_db in 0 s", output)
}
