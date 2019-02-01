package db

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

//todo create more tests
func TestOpenConnection(t *testing.T) {
	sql.Register("fake_driver", &FakeSqlDriver{})

	_, err := NewDb("conn_str", "fake_driver", 1)

	assert.NoError(t, err)
}
