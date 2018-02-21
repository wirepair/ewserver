package store

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/wirepair/ewserver/types"
)

func TestBoltStore_Open(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	d := NewBoltStore()
	if err := d.Open(&Config{ConnectionString: dbFileName}); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	testCloseDb(d, t)
}

func TestBoltStore_StoreUser(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	d := NewBoltStore()
	if err := d.Open(&Config{ConnectionString: dbFileName}); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	defer testCloseDb(d, t)

	u := testUser("user1")

	if err := d.StoreUser(u); err != nil {
		t.Fatalf("error storing user: %s\n", err)
	}
}

func TestBoltStore_FindUserByName(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	d := NewBoltStore()
	if err := d.Open(&Config{ConnectionString: dbFileName}); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	defer testCloseDb(d, t)

	u := testUser("user1")

	if err := d.StoreUser(u); err != nil {
		t.Fatalf("error storing user: %s\n", err)
	}

	foundUser, err := d.FindUserByUserName(u.UserName)
	if err != nil {
		t.Fatalf("error finding user by username: %s\n", err)
	}

	if bytes.Compare(u.ID, foundUser.ID) != 0 {
		t.Fatalf("user IDs do not match: %s and %s\n", u.ID, foundUser.ID)
	}

	if u.FirstName != foundUser.FirstName || u.LastName != foundUser.LastName {
		t.Fatalf("name fields do not match: %s %s and %s %s\n", u.FirstName, u.LastName, foundUser.FirstName, foundUser.LastName)
	}

	if _, err := d.FindUserByUserName("bleh"); err != nil {
		t.Fatalf("got error finding a user who does not exist: %s\n", err)
	}
}

func TestBoltStore_FindUserByID(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	d := NewBoltStore()
	if err := d.Open(&Config{ConnectionString: dbFileName}); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	defer testCloseDb(d, t)

	u := testUser("user1")

	if err := d.StoreUser(u); err != nil {
		t.Fatalf("error storing user: %s\n", err)
	}

	foundUser, err := d.FindUserByID(u.ID)
	if err != nil {
		t.Fatalf("error finding user by username: %s\n", err)
	}

	if u.UserName != foundUser.UserName {
		t.Fatalf("user names do not match: %s and %s\n", u.ID, foundUser.ID)
	}

	if u.FirstName != foundUser.FirstName || u.LastName != foundUser.LastName {
		t.Fatalf("name fields do not match: %s %s and %s %s\n", u.FirstName, u.LastName, foundUser.FirstName, foundUser.LastName)
	}

	if _, err := d.FindUserByID([]byte("does not exist")); err != nil {
		t.Fatalf("got error finding a user who does not exist: %s\n", err)
	}
}

func TestBoltStore_FindUserByAPIKey(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	d := NewBoltStore()
	if err := d.Open(&Config{ConnectionString: dbFileName}); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	defer testCloseDb(d, t)

	u := testUser("user1")

	u2 := testUser("test2")
	u2.APIKey = "aafdsdfsdfsdfsf1"

	if err := d.StoreUser(u); err != nil {
		t.Fatalf("error storing user: %s\n", err)
	}

	if err := d.StoreUser(u2); err != nil {
		t.Fatalf("error storing user2: %s\n", err)
	}

	foundUser, err := d.FindUserByAPIKey(u.APIKey.Bytes())
	if err != nil {
		t.Fatalf("error finding user by APIKey: %s\n", err)
	}

	if bytes.Compare(u.ID, foundUser.ID) != 0 {
		t.Fatalf("user IDs do not match: %s and %s\n", u.ID, foundUser.ID)
	}

	if u.FirstName != foundUser.FirstName || u.LastName != foundUser.LastName {
		t.Fatalf("name fields do not match: %s %s and %s %s\n", u.FirstName, u.LastName, foundUser.FirstName, foundUser.LastName)
	}

	if _, err := d.FindUserByAPIKey([]byte("blerps")); err != nil {
		t.Fatalf("error finding user who doesn't exist by invalid APIKey: %s\n", err)
	}
}

func TestBoltStore_FindAllUsers(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	d := NewBoltStore()
	if err := d.Open(&Config{ConnectionString: dbFileName}); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	defer testCloseDb(d, t)

	u := testUser("user1")

	u2 := testUser("test2")

	if err := d.StoreUser(u); err != nil {
		t.Fatalf("error storing user: %s\n", err)
	}

	if err := d.StoreUser(u2); err != nil {
		t.Fatalf("error storing user2: %s\n", err)
	}

	foundUsers, err := d.FindAllUsers()
	if err != nil {
		t.Fatalf("error finding users: %s\n", err)
	}

	if len(foundUsers) != 2 {
		t.Fatalf("expected 2 users got: %d\n", len(foundUsers))
	}
}

func TestBoltStore_DeleteUserByName(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	d := NewBoltStore()
	if err := d.Open(&Config{ConnectionString: dbFileName}); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	defer testCloseDb(d, t)

	u := testUser("user1")

	if err := d.StoreUser(u); err != nil {
		t.Fatalf("error storing user: %s\n", err)
	}

	foundUser, err := d.FindUserByUserName(u.UserName)
	if err != nil {
		t.Fatalf("error finding user by username: %s\n", err)
	}

	if err := d.DeleteUserByName(foundUser.UserName); err != nil {
		t.Fatalf("error deleting user by name: %s\n", err)
	}

	if err := d.DeleteUserByName(foundUser.UserName); err != nil {
		t.Fatalf("got error (when shouldn't) when deleting a user who no longer exists")
	}
}

func testUser(userName types.UserName) *types.User {
	u := types.NewUser()
	u.ID = []byte("123")
	u.UserName = userName
	u.FirstName = "test"
	u.LastName = "user"
	u.APIKey = "asdfasdf"
	return u
}

func testRemoveDbFile(dbFileName string, t *testing.T) {
	if err := os.Remove(dbFileName); err != nil {
		t.Fatalf("error removing file: %s\n", err)
	}
}

func testCloseDb(d *BoltStore, t *testing.T) {
	if err := d.Close(); err != nil {
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
