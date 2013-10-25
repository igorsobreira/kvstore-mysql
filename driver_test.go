package kvstoremysql

import (
	"database/sql"
	"testing"

	"github.com/igorsobreira/kvstore/testutil"
)

const info = "root:@tcp(localhost:3306)/kvstore_test?timeout=1s"

func TestRequiredAPI(t *testing.T) {
	testutil.TestRequiredAPI(t, teardown, "mysql", info)
}

func teardown() {
	db, err := sql.Open("mysql", info)
	if err != nil {
		println("failed to teardown: " + err.Error())
	}
	_, err = db.Exec("DROP TABLE kvstore")
	if err != nil {
		println("failed to drop table kvstore: " + err.Error())
	}
}
