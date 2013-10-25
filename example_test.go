package kvstoremysql_test

import (
	"database/sql"
	"fmt"

	"github.com/igorsobreira/kvstore"
	_ "github.com/igorsobreira/kvstore-mysql"
)

const info = "root:@tcp(localhost:3306)/kvstore_example?timeout=1s"

func Example() {
	defer Teardown()

	store, err := kvstore.New("mysql", info)
	if err != nil {
		panic(err)
	}

	err = store.Set("foo", []byte("bar"))
	if err != nil {
		panic(err)
	}

	val, err := store.Get("foo")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(val))
	// Output:
	// bar
}

func Teardown() {
	db, err := sql.Open("mysql", info)
	if err != nil {
		panic("failed to teardown: " + err.Error())
	}
	_, err = db.Exec("DROP TABLE kvstore")
	if err != nil {
		panic("failed to drop table kvstore: " + err.Error())
	}
}
