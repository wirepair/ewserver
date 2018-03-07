package boltdb_test

import (
	"fmt"
	"testing"

	"github.com/wirepair/ewserver/ewserver"
	"github.com/wirepair/ewserver/store/boltdb"
)

const (
	testAPIKey      = "dedb33f"
	testLastAddress = "127.0.0.1"
)

func TestAPIUserService_Create(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	service := boltdb.NewAPIUserService(db.DB())
	if err := service.Init(); err != nil {
		t.Fatalf("error initializing user service: %s\n", err)
	}

	testCreateAPIUser(service, t)

	// attempt to create the same user twice
	u := ewserver.NewAPIUser()
	u.Key = testAPIKey
	if err != nil {
		t.Fatalf("error generating api key: %s\n", err)
	}

	if err := service.Create(u); err != ewserver.ErrUserAlreadyExists {
		t.Fatalf("error should have got already exists error, got: %s\n", err)
	}
}

func TestAPIUserService_APIUser(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	service := boltdb.NewAPIUserService(db.DB())
	if err := service.Init(); err != nil {
		t.Fatalf("error initializing user service: %s\n", err)
	}

	testCreateAPIUser(service, t)

	apiUser, err := service.APIUser(testAPIKey)
	if err != nil {
		t.Fatalf("error getting user: %s\n", err)
	}

	if apiUser.LastAddress != testLastAddress {
		t.Fatalf("LastAddress not match: %s and %s\n", apiUser.LastAddress, testLastAddress)
	}

	_, err = service.APIUser(ewserver.APIKey("nobeef"))
	if err == nil {
		t.Fatalf("error should have got error when user does not exist")
	}

	if err != ewserver.ErrUserNotFound {
		t.Fatalf("error should have been not found, got: %s\n", err)
	}
}

func TestAPIUserService_APIUsers(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	service := boltdb.NewAPIUserService(db.DB())
	if err := service.Init(); err != nil {
		t.Fatalf("error initializing API user service: %s\n", err)
	}

	users, err := service.APIUsers()
	if err != nil {
		t.Fatalf("got error when API users is empty: %s\n", err)
	}

	if len(users) != 0 {
		t.Fatalf("API users should be empty")
	}

	testCreateAPIUsers(service, t)

	users, err = service.APIUsers()
	if err != nil {
		t.Fatalf("error getting API users: %s\n", err)
	}

	if len(users) != 10 {
		t.Fatalf("expected 10 API users got: %d\n", len(users))
	}
}

func TestAPIUserService_Delete(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	service := boltdb.NewAPIUserService(db.DB())
	if err := service.Init(); err != nil {
		t.Fatalf("error initializing API user service: %s\n", err)
	}

	if err := service.Delete(testUserName); err != nil {
		t.Fatalf("error should not have got any errors when deleting non-existent API user: %s\n", err)
	}

	testCreateAPIUser(service, t)

	_, err = service.APIUser(testAPIKey)
	if err != nil {
		t.Fatalf("error getting existing user: %s\n", err)
	}

	if err := service.Delete(testAPIKey); err != nil {
		t.Fatalf("error should not have got any errors when deleting an existing API user: %s\n", err)
	}

	if _, err = service.APIUser(testAPIKey); err != nil && err != ewserver.ErrUserNotFound {
		t.Fatalf("error did not get API user not found error, got: %s\n", err)
	}

}

func testCreateAPIUser(service *boltdb.APIUserService, t *testing.T) {
	var err error

	u := ewserver.NewAPIUser()
	u.Key = testAPIKey
	u.LastAddress = testLastAddress

	if service.Create(u); err != nil {
		t.Fatalf("error generating api key: %s\n", err)
	}
}

func testCreateAPIUsers(service *boltdb.APIUserService, t *testing.T) {
	for i := 0; i < 10; i++ {
		u := ewserver.NewAPIUser()
		u.Key = ewserver.APIKey(fmt.Sprintf("%s%d", testAPIKey, i))
		u.LastAddress = testLastAddress

		if err := service.Create(u); err != nil {
			t.Fatalf("error creating user: %s\n", err)
		}
	}
}
