package database_test

import (
	"database/sql/driver"
	"os"
	"testing"
)

type fakeSQLDriver struct {
	connCalls int
	conn      *fakeSQLConn
	err       error
}

func (d *fakeSQLDriver) Open(name string) (driver.Conn, error) {
	d.connCalls++
	if d.err != nil {
		return nil, d.err
	}
	if d.conn == nil {
		d.conn = &fakeSQLConn{}
	}
	return d.conn, nil
}

func (d *fakeSQLDriver) Reset() {
	d.connCalls = 0
	d.err = nil
	d.conn = &fakeSQLConn{}
}

func (d *fakeSQLDriver) MockError(err error) {
	d.err = err
}

type fakeSQLConn struct {
	closeErr error
	stmtErr  error
	txErr    error
}

func (c *fakeSQLConn) Prepare(query string) (driver.Stmt, error) {
	if c.stmtErr != nil {
		return nil, c.stmtErr
	}
	return nil, nil
}

func (c *fakeSQLConn) Close() error {
	if c.closeErr != nil {
		return c.closeErr
	}
	return nil
}

func (c *fakeSQLConn) Begin() (driver.Tx, error) {
	if c.txErr != nil {
		return nil, c.txErr
	}
	return nil, nil
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
