package db

import (
	"database/sql/driver"
	"io"
)

type RowsData struct {
	Cols     []string
	CurRow   int
	Data     [][]interface{}
	IsClosed bool
}

func (fr *RowsData) Columns() []string {
	return fr.Cols
}

func (fr *RowsData) Close() error {
	fr.IsClosed = true
	return nil
}

func (fr *RowsData) Next(dest []driver.Value) error {
	if fr.CurRow > len(fr.Data)-1 {
		return io.EOF
	}

	for i, rowVal := range fr.Data[fr.CurRow] {
		dest[i] = driver.Value(rowVal)
	}
	fr.CurRow++

	return nil
}

type FakeStmt struct {
	NumInputCount int
	RowsAffected  int64
	RowsSlice     *RowsData
	Args          []driver.Value
}

func (fs *FakeStmt) Close() error {
	return nil
}

func (fs *FakeStmt) NumInput() int {
	return fs.NumInputCount
}

func (fs *FakeStmt) Exec(args []driver.Value) (driver.Result, error) {
	fs.Args = args
	return driver.RowsAffected(fs.RowsAffected), nil
}

func (fs *FakeStmt) Query(args []driver.Value) (driver.Rows, error) {
	fs.Args = args
	return fs.RowsSlice, nil
}

type FakeConn struct {
	queries  []string
	FakeStmt *FakeStmt
	FakeFx   *FakeFx
	IsClosed bool
	Error    error
}

func (fc *FakeConn) Prepare(query string) (driver.Stmt, error) {
	fc.queries = append(fc.queries, query)
	return fc.FakeStmt, fc.Error
}

func (fc *FakeConn) Close() error {
	fc.IsClosed = true
	return nil
}

func (fc *FakeConn) Begin() (driver.Tx, error) {
	fc.FakeFx.IsBegun = true
	return fc.FakeFx, nil
}

type FakeFx struct {
	IsBegun    bool
	IsCommit   bool
	IsRollback bool
}

func (fx *FakeFx) Commit() error {
	fx.IsCommit = true
	return nil
}

func (fx *FakeFx) Rollback() error {
	fx.IsRollback = true
	return nil
}

type FakeSqlDriver struct {
	ConnStr         string
	ErrFuncProvider func() error
	Conn            *FakeConn
}

func NewFakeSqlDriver() *FakeSqlDriver {
	rows := &RowsData{
		Cols:     []string{},
		CurRow:   0,
		Data:     [][]interface{}{},
		IsClosed: false,
	}

	stmt := &FakeStmt{
		0,
		0,
		rows,
		[]driver.Value{},
	}

	conn := &FakeConn{
		[]string{},
		stmt,
		&FakeFx{false, false, false},
		false,
		nil,
	}
	drvr := &FakeSqlDriver{
		ConnStr: "",
		ErrFuncProvider: func() error {
			return nil
		},
		Conn: conn,
	}

	return drvr
}

func (fsq *FakeSqlDriver) Open(connStr string) (driver.Conn, error) {
	fsq.ConnStr = connStr

	return fsq.Conn, fsq.ErrFuncProvider()
}
