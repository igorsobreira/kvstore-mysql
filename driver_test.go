package kvstoremysql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/igorsobreira/kvstore"
)

const (
	info     = "root:@tcp(localhost:3306)/kvstore_mysql_test?timeout=1s"
	Megabyte = 1024 * 1024
)

// Basic tests that sets a value, gets it then delete it
//
// Tests the limits of the keys and values
func TestSetGetDelete(t *testing.T) {
	defer Teardown()

	store, err := kvstore.New("mysql", info)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		Key string
		Val []byte
	}{
		{
			Key: "key1",
			Val: []byte("value1"),
		},
		{
			Key: "key2",
			Val: ByteSlice('V', 1*Megabyte), // max value size by default (max_long_data_size or max_allowed_packet mysql options)
		},
		{
			Key: String('K', 256), // max key size
			Val: []byte("value3"),
		},
	}

	for _, tt := range tests {
		if err := store.Set(tt.Key, tt.Val); err != nil {
			t.Errorf("set %#v failed: %s", tt.Key, err)
			continue
		}

		val, err := store.Get(tt.Key)
		if err != nil {
			t.Errorf("get %#v failed: %s", tt.Key, err)
			continue
		}
		if !ByteSliceEqual(val, tt.Val) {
			t.Errorf("get %#v got %#v, want %#v", tt.Key, Truncate(val), Truncate(tt.Val))
			continue
		}

		if err = store.Delete(tt.Key); err != nil {
			t.Errorf("delete %#v failed: %s", tt.Key, err)
			continue
		}

		val, err = store.Get(tt.Key)
		if err != kvstore.ErrNotFound {
			t.Errorf("invalid error for key (%#v) not found: %s", tt.Key, err)
		}
		if val != nil {
			t.Errorf("get %#v after delete should return nil, found %#v", tt.Key, Truncate(val))
		}
	}
}

func TestSetOverride(t *testing.T) {
	defer Teardown()

	store, err := kvstore.New("mysql", info)
	if err != nil {
		t.Fatal(err)
	}

	err = store.Set("key", []byte("value1"))
	if err != nil {
		t.Error("first set:", err)
	}
	err = store.Set("key", []byte("value2"))
	if err != nil {
		t.Error("second set:", err)
	}

	val, err := store.Get("key")
	if err != nil {
		t.Error("get:", err)
	}
	if !ByteSliceEqual(val, []byte("value2")) {
		t.Error("got:", string(val))
	}
}

func TestGetNotFound(t *testing.T) {
	defer Teardown()

	store, err := kvstore.New("mysql", info)
	if err != nil {
		t.Fatal(err)
	}

	val, err := store.Get("key")

	if err != kvstore.ErrNotFound {
		t.Error("invalid error:", err)
	}
	if val != nil {
		t.Errorf("got: %#v, want nil", val)
	}
}

func TestDeleteNotFound(t *testing.T) {
	defer Teardown()

	store, err := kvstore.New("mysql", info)
	if err != nil {
		t.Fatal(err)
	}

	err = store.Delete("something")

	if err != nil {
		t.Error(err)
	}
}

// helpers

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

func ByteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func String(char byte, repeat int) string {
	return string(ByteSlice(char, repeat))
}

func ByteSlice(char byte, repeat int) []byte {
	result := make([]byte, repeat)
	for i := 0; i < repeat; i++ {
		result[i] = char
	}
	return result
}

func Truncate(s []byte) string {
	max := 10
	if len(s) <= max {
		return fmt.Sprintf("%#v", s)
	}
	return fmt.Sprintf("%#v (truncated, size %d)", s[:max], len(s))
}
