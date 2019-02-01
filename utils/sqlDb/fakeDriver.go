package db

import "database/sql/driver"

type FakeRows struct {}

func (fr *FakeRows) Columns() []string {
	return []string{}
}

func (fr *FakeRows) Close() error {
	return nil
}

func (fr *FakeRows) Next(dest []driver.Value) error {
	return nil
}

type FakeStmt struct {}

func (fs *FakeStmt) Close() error {
	return nil
}

func (fs *FakeStmt) NumInput() int {
	return 0
}

func (fs *FakeStmt) Exec(args []driver.Value) (driver.Result, error) {
	return driver.RowsAffected(0), nil
}

func (fs *FakeStmt) Query(args []driver.Value) (driver.Rows, error) {
	return &FakeRows{}, nil
}

type FakeConn struct {

}

func (fc *FakeConn) Prepare(query string) (driver.Stmt, error) {
	return &FakeStmt{}, nil
}

func (fc *FakeConn) Close() error {
	return nil
}

func (fc *FakeConn) Begin() (driver.Tx, error) {
	return &FakeFx{}, nil
}

type FakeFx struct {

}

func (fx *FakeFx) Commit() error {
	return nil
}

func (fx *FakeFx) Rollback() error {
	return nil
}

type FakeSqlDriver struct {}

func (fsq *FakeSqlDriver) Open(name string) (driver.Conn, error) {
	return &FakeConn{}, nil
}
