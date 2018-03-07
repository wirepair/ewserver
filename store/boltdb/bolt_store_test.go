package boltdb_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/wirepair/ewserver/store"
	"github.com/wirepair/ewserver/store/boltdb"
)

func TestBoltStore_Open(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)
	db := testOpenDb(dbFileName, t)
	testCloseDb(db, t)
}

func testRemoveDbFile(dbFileName string, t *testing.T) {
	if err := os.Remove(dbFileName); err != nil {
		t.Fatalf("error removing file: %s\n", err)
	}
}

func testOpenDb(dbFileName string, t *testing.T) *boltdb.BoltStore {
	config := store.NewConfig()
	config.Options["database"] = dbFileName

	db := boltdb.NewBoltStore()
	if err := db.Open(config); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	return db
}

func testCloseDb(db *boltdb.BoltStore, t *testing.T) {
	if err := db.Close(); err != nil {
		t.Fatalf("error closing database: %s\n", err)
	}
}

func testTempDbFileName(dir string) (string, error) {
	f, err := ioutil.TempFile(dir, "db")
	if err != nil {
		return "", err
	}

	f.Close()
	os.Remove(f.Name())

	return f.Name(), nil
}
