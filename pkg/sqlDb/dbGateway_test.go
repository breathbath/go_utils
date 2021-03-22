package db

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initDBGateway(driverID string) (*Gateway, *FakeSQLDriver, error) {
	fakeDriver := NewFakeSQLDriver()

	sql.Register(driverID, fakeDriver)

	sqlX, err := NewDB("conn_str", driverID)

	if err != nil {
		return nil, fakeDriver, err
	}

	dbgatewy := NewDBGateway(sqlX)

	return dbgatewy, fakeDriver, nil
}

type SomeSlice struct {
	Number int    `db:"number"`
	Color  string `db:"color"`
	Object string `db:"object"`
}

func TestFindByQuery(t *testing.T) {
	dbGateway, fakeDriver, err := initDBGateway("TestFindByQuery_driver")
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

	fakeDriver.Conn.Error = errors.New("wrong query syntax")
	err = dbGateway.FindByQuery(&data, "Badsql", []interface{}{2})
	assert.EqualError(t, err, "query 'Badsql' with args [2] has failed: wrong query syntax")
}

func TestQueryWithCallback(t *testing.T) {
	dbGateway, fakeDriver, err := initDBGateway("TestQueryWithCallback_driver")
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

	fakeDriver.Conn.Error = errors.New("mismatch")
	err = dbGateway.QueryWithCallback(
		func(row map[string]interface{}, errCallback error) {},
		"badQ",
		map[string]interface{}{},
	)
	assert.EqualError(t, err, "query 'badQ' with args map[] has failed: mismatch")
}

func TestExec(t *testing.T) {
	dbGateway, fakeDriver, err := initDBGateway("TestExec_driver")
	assert.NoError(t, err)

	_, err = dbGateway.Exec("Some exec query")
	assert.NoError(t, err)
	assert.Equal(t, []string{"Some exec query"}, fakeDriver.Conn.queries)

	fakeDriver.Conn.Error = errors.New("wrong query syntax")
	_, err = dbGateway.Exec("Bad query", 1, 2)
	assert.EqualError(t, err, "query 'Bad query' with args [1 2] has failed: wrong query syntax")
}

func TestTransactions(t *testing.T) {
	dbGateway, fakeDriver, err := initDBGateway("TestTransactions_driver")
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
	dbGateway, fakeDriver, err := initDBGateway("TestTruncate_driver")
	assert.NoError(t, err)

	_, err = dbGateway.TruncateTable("mytable")
	assert.NoError(t, err)
	assert.Equal(t, []string{"TRUNCATE TABLE `mytable`"}, fakeDriver.Conn.queries)
}

func TestDestroy(t *testing.T) {
	dbGateway, fakeDriver, err := initDBGateway("TestDestroy_driver")
	assert.NoError(t, err)

	_, err = dbGateway.Exec("Some query to open conn")
	assert.NoError(t, err)
	if err != nil {
		return
	}

	err = dbGateway.Destroy()
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.True(t, fakeDriver.Conn.IsClosed)
}

func TestFindByQueryFlex(t *testing.T) {
	dbGateway, fakeDriver, err := initDBGateway("TestFindByQueryFlex_driver")
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

	fakeDriver.Conn.Error = errors.New("bad arguments count")
	err = dbGateway.FindByQueryFlex(&data, "Badsql")
	assert.EqualError(t, err, "query 'Badsql' with args [] has failed: bad arguments count")
}

func TestScanScalarByQuery(t *testing.T) {
	dbGateway, fakeDriver, err := initDBGateway("TestScanScalarByQuery_driver")
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

	fakeDriver.Conn.Error = errors.New("wrong scalar type")
	_, err = dbGateway.ScanScalarByQuery(&numb, "Sql")
	assert.EqualError(t, err, "query 'Sql' with args [] has failed: wrong scalar type")

	fakeDriver.Conn.Error = sql.ErrNoRows
	found, err = dbGateway.ScanScalarByQuery(&numb, "Sql")
	assert.False(t, found)
	assert.NoError(t, err)
}

func TestScanStructByQuery(t *testing.T) {
	dbGateway, fakeDriver, err := initDBGateway("TestScanStructByQuery_driver")
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

	fakeDriver.Conn.Error = errors.New("wrong struct field")
	_, err = dbGateway.ScanStructByQuery(&data, "Sql")
	assert.EqualError(t, err, "query 'Sql' with args [] has failed: wrong struct field")

	fakeDriver.Conn.Error = sql.ErrNoRows
	found, err = dbGateway.ScanStructByQuery(&data, "Sql")
	assert.False(t, found)
	assert.NoError(t, err)
}

func TestFindOneStructById(t *testing.T) {
	dbGateway, fakeDriver, err := initDBGateway("TestFindOneStructById_driver")
	assert.NoError(t, err)

	fakeDriver.Conn.FakeStmt.NumInputCount = 1

	fakeDriver.Conn.FakeStmt.RowsSlice.Data = [][]interface{}{
		{222, "red", "baloon"},
	}
	fakeDriver.Conn.FakeStmt.RowsSlice.Cols = []string{"number", "color", "object"}

	data := SomeSlice{}

	found, err := dbGateway.FindOneStructByID(
		&data,
		"objects_table",
		333,
	)
	assert.NoError(t, err)
	assert.True(t, found)

	assert.Equal(t, []string{"Select `objects_table`.* from `objects_table` WHERE `objects_table`.id=? LIMIT 1"}, fakeDriver.Conn.queries)

	assert.Equal(
		t,
		SomeSlice{222, "red", "baloon"},
		data,
	)

	assert.Equal(t, []driver.Value{int64(333)}, fakeDriver.Conn.FakeStmt.Args)
}
