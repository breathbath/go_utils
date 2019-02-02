package db

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initDbGateway(driverId string) (*DbGateway, *FakeSqlDriver, error) {
	fakeDriver := NewFakeSqlDriver()

	sql.Register(driverId, fakeDriver)

	sqlX, err := NewDb("conn_str", driverId, 1)

	if err != nil {
		return nil, fakeDriver, err
	}

	dbgatewy := NewDbGateway(sqlX)

	return dbgatewy, fakeDriver, nil
}

type SomeSlice struct {
	Number int    `db:"number"`
	Color  string `db:"color"`
	Object string `db:"object"`
}

func TestFindByQuery(t *testing.T) {
	dbGateway, fakeDriver, err := initDbGateway("TestFindByQuery_driver")
	assert.NoError(t, err)

	fakeDriver.Conn.FakeStmt.RowsSlice.Data = [][]interface{}{
		{1, "red", "ball"},
		{2, "green", "pen"},
	}
	fakeDriver.Conn.FakeStmt.RowsSlice.Cols = []string{"number", "color", "object"}

	data := []SomeSlice{}

	err = dbGateway.FindByQuery(
		&data,
		"some sql query",
		[]interface{}{},
	)
	assert.NoError(t, err)

	assert.Equal(t, []string{"some sql query"}, fakeDriver.Conn.queries)

	assert.Equal(
		t,
		[]SomeSlice{
			{1, "red", "ball"},
			{2, "green", "pen"},
		},
		data,
	)

	fakeDriver.Conn.Error = errors.New("Wrong query syntax")
	err = dbGateway.FindByQuery(&data, "Badsql", []interface{}{2})
	assert.EqualError(t, err, "Query 'Badsql' with args [2] has failed: Wrong query syntax")
}

func TestQueryWithCallback(t *testing.T) {
	dbGateway, fakeDriver, err := initDbGateway("TestQueryWithCallback_driver")
	assert.NoError(t, err)

	fakeDriver.Conn.FakeStmt.RowsSlice.Data = [][]interface{}{
		{1, "red", "ball"},
		{2, "green", "ball"},
	}
	fakeDriver.Conn.FakeStmt.RowsSlice.Cols = []string{"number", "color", "object"}

	receivedResult := []map[string]interface{}{}
	err = dbGateway.QueryWithCallback(
		func(row map[string]interface{}, errCallback error) {
			receivedResult = append(receivedResult, row)
		},
		"some query",
		map[string]interface{}{},
	)
	assert.NoError(t, err)

	assert.Equal(t, []string{"some query"}, fakeDriver.Conn.queries)
	assert.Equal(
		t,
		[]map[string]interface{}{
			{"number": 1, "color": "red", "object": "ball"},
			{"number": 2, "color": "green", "object": "ball"},
		},
		receivedResult,
	)
	assert.True(t, fakeDriver.Conn.FakeStmt.RowsSlice.IsClosed)

	fakeDriver.Conn.Error = errors.New("Mismatch")
	err = dbGateway.QueryWithCallback(
		func(row map[string]interface{}, errCallback error) {},
		"badQ",
		map[string]interface{}{},
	)
	assert.EqualError(t, err, "Query 'badQ' with args map[] has failed: Mismatch")
}

func TestExec(t *testing.T) {
	dbGateway, fakeDriver, err := initDbGateway("TestExec_driver")
	assert.NoError(t, err)

	_, err = dbGateway.Exec("Some exec query")
	assert.NoError(t, err)
	assert.Equal(t, []string{"Some exec query"}, fakeDriver.Conn.queries)

	fakeDriver.Conn.Error = errors.New("Wrong query syntax")
	_, err = dbGateway.Exec("Bad query", 1, 2)
	assert.EqualError(t, err, "Query 'Bad query' with args [1 2] has failed: Wrong query syntax")
}

func TestTransactions(t *testing.T) {
	dbGateway, fakeDriver, err := initDbGateway("TestTransactions_driver")
	assert.NoError(t, err)

	lastFx := fakeDriver.Conn.FakeFx
	assert.False(t, lastFx.IsBegun || lastFx.IsRollback || lastFx.IsCommit)

	tx, err := dbGateway.Begin()
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)

	assert.True(t, lastFx.IsBegun && !lastFx.IsRollback && lastFx.IsCommit)
	lastFx.IsCommit = false
	lastFx.IsRollback = false
	lastFx.IsBegun = false

	tx, err = dbGateway.Begin()
	assert.NoError(t, err)

	err = tx.Rollback()
	assert.NoError(t, err)
	assert.True(t, lastFx.IsBegun && lastFx.IsRollback && !lastFx.IsCommit)
}

func TestTruncate(t *testing.T) {
	dbGateway, fakeDriver, err := initDbGateway("TestTruncate_driver")
	assert.NoError(t, err)

	_, err = dbGateway.TruncateTable("mytable")
	assert.NoError(t, err)
	assert.Equal(t, []string{"TRUNCATE TABLE `mytable`"}, fakeDriver.Conn.queries)
}

func TestDestroy(t *testing.T) {
	dbGateway, fakeDriver, err := initDbGateway("TestDestroy_driver")
	assert.NoError(t, err)

	dbGateway.Exec("Some query to open conn")

	err = dbGateway.Destroy()
	assert.NoError(t, err)

	assert.True(t, fakeDriver.Conn.IsClosed)
}

func TestFindByQueryFlex(t *testing.T) {
	dbGateway, fakeDriver, err := initDbGateway("TestFindByQueryFlex_driver")
	assert.NoError(t, err)

	fakeDriver.Conn.FakeStmt.NumInputCount = 2

	fakeDriver.Conn.FakeStmt.RowsSlice.Data = [][]interface{}{
		{222, "yellow", "submarine"},
	}
	fakeDriver.Conn.FakeStmt.RowsSlice.Cols = []string{"number", "color", "object"}

	data := []SomeSlice{}

	err = dbGateway.FindByQueryFlex(
		&data,
		"some sql query",
		1,
		"param2",
	)
	assert.NoError(t, err)

	assert.Equal(t, []string{"some sql query"}, fakeDriver.Conn.queries)

	assert.Equal(
		t,
		[]SomeSlice{
			{222, "yellow", "submarine"},
		},
		data,
	)

	assert.Equal(t, []driver.Value{int64(1), "param2"}, fakeDriver.Conn.FakeStmt.Args)

	fakeDriver.Conn.Error = errors.New("Bad arguments count")
	err = dbGateway.FindByQueryFlex(&data, "Badsql")
	assert.EqualError(t, err, "Query 'Badsql' with args [] has failed: Bad arguments count")
}

func TestScanScalarByQuery(t *testing.T) {
	dbGateway, fakeDriver, err := initDbGateway("TestScanScalarByQuery_driver")
	assert.NoError(t, err)

	fakeDriver.Conn.FakeStmt.RowsSlice.Cols = []string{"number"}
	fakeDriver.Conn.FakeStmt.RowsSlice.Data = [][]interface{}{{1}}

	var numb int64

	found, err := dbGateway.ScanScalarByQuery(
		&numb,
		"some scalar query",
	)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, int64(1), numb)

	fakeDriver.Conn.Error = errors.New("Wrong scalar type")
	_, err = dbGateway.ScanScalarByQuery(&numb, "Sql")
	assert.EqualError(t, err, "Query 'Sql' with args [] has failed: Wrong scalar type")
}

func TestScanStructByQuery(t *testing.T) {
	dbGateway, fakeDriver, err := initDbGateway("TestScanStructByQuery_driver")
	assert.NoError(t, err)

	fakeDriver.Conn.FakeStmt.RowsSlice.Data = [][]interface{}{
		{222, "red", "baloon"},
	}
	fakeDriver.Conn.FakeStmt.RowsSlice.Cols = []string{"number", "color", "object"}

	data := SomeSlice{}

	found, err := dbGateway.ScanStructByQuery(
		&data,
		"some q",
	)
	assert.NoError(t, err)
	assert.True(t, found)

	assert.Equal(t, []string{"some q"}, fakeDriver.Conn.queries)

	assert.Equal(
		t,
		SomeSlice{222, "red", "baloon"},
		data,
	)

	fakeDriver.Conn.Error = errors.New("Wrong struct field")
	_, err = dbGateway.ScanScalarByQuery(&data, "Sql")
	assert.EqualError(t, err, "Query 'Sql' with args [] has failed: Wrong struct field")
}