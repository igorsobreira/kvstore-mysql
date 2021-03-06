// Package kvstoremysql implements a MySQL driver for github.com/igorsobreira/kvstore.
//
// Keys are stored in a VARCHAR(256) column and values in a MEDIUMBLOB.
// The max size of a value is specified by MySQL config max_long_data_size and
// max_allowed_packet on version 5.6
package kvstoremysql

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"io"

	_ "github.com/go-sql-driver/mysql" // register mysql driver
	"github.com/igorsobreira/kvstore"
)

func init() {
	kvstore.Register("mysql", &Driver{})
}

// Driver implementes kvstore.Driver interface
type Driver struct{}

// Conn implements kvstore.Conn wrapping a MySQL connection
type Conn struct {
	db *sql.DB
}

// Open is called by kvstore.New(), will open the connection to MySQL
//
// info has to be the Data Source Format, as specified in
// https://github.com/Go-SQL-Driver/MySQL/#dsn-data-source-name
//
// The database specified on info has to exist, the necessary tables
// will be created
func (d *Driver) Open(info string) (kvstore.Conn, error) {

	db, err := sql.Open("mysql", info)
	if err != nil {
		return nil, err
	}

	// create a key_md5 as UNIQUE since key is too big to create the index,
	// and I want to use INSERT ... ON DUPLICATE KEY UPDATE ...
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS kvstore (" +
		"`id`      INT          NOT NULL auto_increment PRIMARY KEY," +
		"`key`     VARCHAR(256) NOT NULL," +
		"`key_md5` CHAR(32)     NOT NULL UNIQUE," +
		"`value`   MEDIUMBLOB   NOT NULL" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8")
	if err != nil {
		db.Close()
		return nil, errors.New("kvstore mysql: " + err.Error())
	}

	return &Conn{db}, nil
}

// Set key associated to value. Override existing value.
func (c *Conn) Set(key string, value []byte) error {
	_, err := c.db.Exec(
		"INSERT INTO kvstore (`key`, `key_md5`, `value`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `value` = ?",
		key, md5of(key), value, value,
	)
	return err
}

// Get value associated with key. Return kvstore.ErrNotFound if
// key doesn't exist
func (c *Conn) Get(key string) (value []byte, err error) {
	err = c.db.QueryRow("SELECT `value` FROM kvstore WHERE `key`=?", key).Scan(&value)

	switch {
	case err == sql.ErrNoRows:
		return value, kvstore.ErrNotFound
	case err != nil:
		return value, err
	default:
		return value, nil
	}
}

// Delete key. No-op if key not found.
func (c *Conn) Delete(key string) (err error) {
	_, err = c.db.Exec("DELETE FROM kvstore WHERE `key` = ?", key)
	return err
}

// Close will close the mysql connection.
func (c *Conn) Close() error {
	return c.db.Close()
}

func md5of(s string) string {
	m := md5.New()
	io.WriteString(m, s)
	return fmt.Sprintf("%x", m.Sum(nil))
}
