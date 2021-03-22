package db

import (
	baseSql "database/sql"
	"fmt"

	"github.com/breathbath/go_utils/v2/io"
	"github.com/jmoiron/sqlx"
)

// Gateway wrapper for sqlx with some useful additions
type Gateway struct {
	conn *sqlx.DB
}

func NewDBGateway(conn *sqlx.DB) *Gateway {
	return &Gateway{conn: conn}
}

// QueryWithCallback executes query and gives each row to callback func
func (dg *Gateway) QueryWithCallback(
	resultCallback func(row map[string]interface{}, errCallback error),
	sql string,
	args map[string]interface{},
) error {
	rows, err := dg.conn.NamedQuery(sql, args)
	err = dg.packError(err, sql, args)
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

// Begin transaction
func (dg *Gateway) Begin() (*baseSql.Tx, error) {
	return dg.conn.Begin()
}

// Exec executes a query, very useful for mutations like insert, delete, update, alter, drop, truncate etc.
func (dg *Gateway) Exec(query string, args ...interface{}) (baseSql.Result, error) {
	res, err := dg.conn.Exec(query, args...)
	err = dg.packError(err, query, args)

	return res, err
}

// TruncateTable truncates table
func (dg *Gateway) TruncateTable(tableName string) (baseSql.Result, error) {
	res, err := dg.conn.Exec("TRUNCATE TABLE " + dg.escapeTableName(tableName))

	return res, err
}

// Destroy closes connection, useful for garbage collection
func (dg *Gateway) Destroy() error {
	return dg.conn.Close()
}

// FindByQuery executes a select query, saving the result into target, args are packed into a slice
func (dg *Gateway) FindByQuery(target interface{}, sql string, args []interface{}) error {
	err := dg.conn.Select(target, sql, args...)

	return dg.packError(err, sql, args)
}

// FindByQueryFlex the same as FindByQuery but arguments can be omitted
func (dg *Gateway) FindByQueryFlex(target interface{}, sql string, args ...interface{}) error {
	err := dg.conn.Select(target, sql, args...)

	return dg.packError(err, sql, args)
}

// ScanScalarByQuery useful to get simple scalar results from db
func (dg *Gateway) ScanScalarByQuery(target interface{}, sql string, args ...interface{}) (found bool, err error) {
	err = dg.conn.QueryRow(sql, args...).Scan(target)

	if err == baseSql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, dg.packError(err, sql, args)
	}

	return true, nil
}

func (dg *Gateway) FindOneStructByID(target interface{}, tableName string, id int64) (bool, error) {
	tableName = dg.escapeTableName(tableName)
	q := fmt.Sprintf("Select %s.* from %s WHERE %s.id=? LIMIT 1", tableName, tableName, tableName)

	return dg.ScanStructByQuery(target, q, id)
}

// ScanStructByQuery useful to get a single struct
func (dg *Gateway) ScanStructByQuery(target interface{}, q string, args ...interface{}) (bool, error) {
	err := dg.conn.Get(target, q, args...)

	if err == baseSql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, dg.packError(err, q, args)
	}

	return true, nil
}

func (dg *Gateway) packError(err error, query string, args interface{}) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(
		"query '%s' with args %v has failed: %v",
		io.RemoveLineBreaks(query),
		args,
		err,
	)
}

func (dg *Gateway) escapeTableName(tableName string) string {
	if tableName[0] != '`' {
		tableName = fmt.Sprintf("`%s`", tableName)
	}

	return tableName
}
