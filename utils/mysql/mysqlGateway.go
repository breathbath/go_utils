package db

import (
	"database/sql"
	"fmt"
	"github.com/breathbath/go_utils/utils/io"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Mysql struct {
	conn *sqlx.DB
}

func NewMysql(conn *sqlx.DB) *Mysql {
	return &Mysql{conn: conn}
}

func (m *Mysql) Query(sql string, args map[string]interface{}, resultCallback func(row map[string]interface{}, errCallback error)) error {
	rows, err := m.conn.NamedQuery(sql, args)
	err = m.packError(err, sql, args)
	if err != nil {
		return err
	}

	for rows.Next() {
		result := map[string]interface{}{}
		err = rows.MapScan(result)
		resultCallback(result, err)
	}
	defer rows.Close()

	return nil
}

func (m *Mysql) escapeTableName(tableName string) string {
	if tableName[0] != '`' {
		tableName = fmt.Sprintf("`%s`", tableName)
	}

	return tableName
}

func (m *Mysql) Begin() (*sql.Tx, error) {
	return m.conn.Begin()
}

func (m *Mysql) Exec(query string, args ...interface{}) (sql.Result, error) {
	res, err := m.conn.Exec(query, args...)
	err = m.packError(err, query, args)

	return res, err
}

func (m *Mysql) TruncateTable(tableName string) (sql.Result, error) {
	res := m.conn.MustExec("TRUNCATE TABLE " + m.escapeTableName(tableName))

	return res, nil
}

func (m *Mysql) Destroy() error {
	if m.conn != nil {
		return m.conn.Close()
	}

	return nil
}

func (m *Mysql) FindByQuery(sql string, args []interface{}, obj interface{}) error {
	err := m.conn.Select(obj, sql, args...)

	return m.packError(err, sql, args)
}

func (m *Mysql) FindBy(sql string, target interface{}, args ...interface{}) error {
	err := m.conn.Select(target, sql, args...)

	return m.packError(err, sql, args)
}

func (m *Mysql) ScanScalarByQuery(q string, args []interface{}, obj interface{}) (bool, error) {
	err := m.conn.QueryRow(q, args...).Scan(obj)

	if err == sql.ErrNoRows {
		return false, nil
	}

	err = m.packError(err, q, args)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *Mysql) ScanStructByQuery(q string, args []interface{}, obj interface{}) (bool, error) {
	err := m.conn.Get(obj, q, args...)

	if err == sql.ErrNoRows {
		return false, nil
	}

	err = m.packError(err, q, args)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *Mysql) packError(err error, query string, args interface{}) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(
		"Query %s with args %v has failed: %v",
		io.RemoveLineBreaks(query),
		args,
		err,
	)
}
