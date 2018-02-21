package repositories

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/wirepair/ewserver/store"
	"github.com/wirepair/ewserver/types"
)

func TestUserRepositoryAdd(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	d := store.NewBoltStore()
	if err := d.Open(&store.Config{ConnectionString: dbFileName}); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	defer testCloseDb(d, t)

	userRepo := NewUserRepository(d)

	u := testUser(types.UserName("user1"))
	if err := userRepo.AddUser(u); err != nil {
		t.Fatalf("error adding user to userrepo: %s\n", err)
	}

	// test adding user that already exists
	if err := userRepo.AddUser(u); err == nil {
		t.Fatalf("did not get error adding user that already exists\n")
	}
}

func TestUserRepositoryDelete(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	d := store.NewBoltStore()
	if err := d.Open(&store.Config{ConnectionString: dbFileName}); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	defer testCloseDb(d, t)

	userRepo := NewUserRepository(d)

	u := testUser(types.UserName("user1"))
	if err := userRepo.DeleteUser(u); err == nil {
		t.Fatalf("did not get error when deleting a user who does not exist\n")
	}

	if err := userRepo.AddUser(u); err != nil {
		t.Fatalf("error adding user to userrepo: %s\n", err)
	}

	if err := userRepo.DeleteUser(u); err != nil {
		t.Fatalf("got error when attempting to delete a user who exists: %s\n", err)
	}

	if err := userRepo.AddUser(u); err != nil {
		t.Fatalf("error adding user to userrepo: %s\n", err)
	}

	foundUser, err := userRepo.FindByUserName(u.UserName)
	if err != nil {
		t.Fatalf("error finding user who should exist: %s\n", err)
	}

	if foundUser == nil {
		t.Fatalf("unable to find user who should exist\n")
	}
}

func TestUserRepositoryUpdateUser(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	d := store.NewBoltStore()
	if err := d.Open(&store.Config{ConnectionString: dbFileName}); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	defer testCloseDb(d, t)

	userRepo := NewUserRepository(d)

	u := testUser("user1")

	if err := userRepo.AddUser(u); err != nil {
		t.Fatalf("error adding user to userrepo: %s\n", err)
	}

	updatedUser := testUser("user1")
	updatedUser.FirstName = "updated"
	updatedUser.LastName = "updated"
	updatedUser.APIKey = "updated"

	if err := userRepo.UpdateUser(updatedUser); err != nil {
		t.Fatalf("error updating user: %s\n", err)
	}

	foundUser, err := userRepo.FindByUserName(u.UserName)
	if err != nil {
		t.Fatalf("error finding user who should exist: %s\n", err)
	}

	if foundUser == nil {
		t.Fatalf("unable to find user who should exist\n")
	}

	if foundUser.LastName != "updated" || foundUser.FirstName != "updated" || foundUser.APIKey != "updated" {
		t.Fatalf("store did not reflect updated changes: %#v\n", foundUser)
	}

	nonexist := testUser("user2")
	if err := userRepo.UpdateUser(nonexist); err == nil {
		t.Fatalf("did not get error updating non-existent user\n")
	}
}

func TestUserRepositoryFindUser(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	d := store.NewBoltStore()
	if err := d.Open(&store.Config{ConnectionString: dbFileName}); err != nil {
		t.Fatalf("error opening database file: %s\n", err)
	}
	defer testCloseDb(d, t)

	userRepo := NewUserRepository(d)

	u := testUser("user1")
	u2 := testUser("user2")
	u2.APIKey = "apikey2"

	if err := userRepo.AddUser(u); err != nil {
		t.Fatalf("error adding user to userrepo: %s\n", err)
	}

	if err := userRepo.AddUser(u2); err != nil {
		t.Fatalf("error adding user to userrepo: %s\n", err)
	}

	foundu, err := userRepo.FindByAPIKey(u.APIKey)
	if err != nil {
		t.Fatalf("error finding by apikey: %s\n", err)
	}

	if foundu.UserName != u.UserName {
		t.Fatalf("found user name by apikey did not match original: %s != %s\n", foundu.UserName, u.UserName)
	}

	foundu, err = userRepo.FindByUserName(u.UserName)
	if err != nil {
		t.Fatalf("error finding by apikey: %s\n", err)
	}

	if foundu.UserName != u.UserName {
		t.Fatalf("found user name by username did not match original: %s != %s\n", foundu.UserName, u.UserName)
	}

	foundu, err = userRepo.FindByID(u.ID)
	if err != nil {
		t.Fatalf("error finding by apikey: %s\n", err)
	}

	if foundu.UserName != u.UserName {
		t.Fatalf("found user name by username did not match original: %s != %s\n", foundu.UserName, u.UserName)
	}

	foundUsers, err := userRepo.FindAllUsers()
	if err != nil {
		t.Fatalf("error finding all users: %s\n", err)
	}

	if len(foundUsers) != 2 {
		t.Fatalf("expected 2 users, got: %d\n", len(foundUsers))
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

func testCloseDb(d store.Storer, t *testing.T) {
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
