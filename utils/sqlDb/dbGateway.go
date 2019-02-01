package db

import (
	baseSql "database/sql"
	"fmt"
	"github.com/breathbath/go_utils/utils/io"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//DbGateway wrapper for sqlx with some useful additions
type DbGateway struct {
	conn *sqlx.DB
}

func NewDbGateway(conn *sqlx.DB) *DbGateway {
	return &DbGateway{conn: conn}
}

//QueryWithCallback eecutes query and gives each row to callback func
func (dg *DbGateway) QueryWithCallback(
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

//Begin transaction
func (dg *DbGateway) Begin() (*baseSql.Tx, error) {
	return dg.conn.Begin()
}

//Exec executes a query, very useful for mutations like insert, delete, update, alter, drop, truncate etc.
func (dg *DbGateway) Exec(query string, args ...interface{}) (baseSql.Result, error) {
	res, err := dg.conn.Exec(query, args...)
	err = dg.packError(err, query, args)

	return res, err
}

//TruncateTable truncates table
func (dg *DbGateway) TruncateTable(tableName string) (baseSql.Result, error) {
	res, err := dg.conn.Exec("TRUNCATE TABLE " + dg.escapeTableName(tableName))

	return res, err
}

//Destroy closes connection, useful for garbage collection
func (dg *DbGateway) Destroy() error {
	return dg.conn.Close()
}

//FindByQuery executes a select query, saving the result into target, args are packed into a slice
func (dg *DbGateway) FindByQuery(target interface{}, sql string, args []interface{}) error {
	err := dg.conn.Select(target, sql, args...)

	return dg.packError(err, sql, args)
}

//FindByQueryFlex the same as FindByQuery but arguments can be omitted
func (dg *DbGateway) FindByQueryFlex(target interface{}, sql string, args ...interface{}) error {
	err := dg.conn.Select(target, sql, args...)

	return dg.packError(err, sql, args)
}

//ScanScalarByQuery useful to get simple scalar results from db
func (dg *DbGateway) ScanScalarByQuery(target interface{}, sql string, args ...interface{}) (found bool, err error) {
	err = dg.conn.QueryRow(sql, args...).Scan(target)

	if err == baseSql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, dg.packError(err, sql, args)
	}

	return true, nil
}

//ScanStructByQuery useful to get a single struct
func (dg *DbGateway) ScanStructByQuery(q string, args []interface{}, obj interface{}) (bool, error) {
	err := dg.conn.Get(obj, q, args...)

	if err == baseSql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, dg.packError(err, q, args)
	}

	return true, nil
}

func (dg *DbGateway) packError(err error, query string, args interface{}) error {
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

func (dg *DbGateway) escapeTableName(tableName string) string {
	if tableName[0] != '`' {
		tableName = fmt.Sprintf("`%s`", tableName)
	}

	return tableName
}
