package boltdb_test

import (
	"fmt"
	"testing"

	"github.com/wirepair/ewserver/ewserver"
	"github.com/wirepair/ewserver/store/boltdb"
)

const (
	testUserName  = "user1"
	testFirstName = "user"
	testLastName  = "1"
	testPassword  = "somepassword"
)

func TestUserService_Create(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	service := boltdb.NewUserService(db.DB())
	if err := service.Init(); err != nil {
		t.Fatalf("error initializing user service: %s\n", err)
	}

	testCreateUser(service, t)

	// attempt to create the same user twice
	u := ewserver.NewUser()
	u.UserName = testUserName
	if err := service.Create(u, "someotherpass"); err != ewserver.ErrUserAlreadyExists {
		t.Fatalf("error should have got already exists error, got: %s\n", err)
	}
}

func TestUserService_User(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	service := boltdb.NewUserService(db.DB())
	if err := service.Init(); err != nil {
		t.Fatalf("error initializing user service: %s\n", err)
	}

	testCreateUser(service, t)

	user, err := service.User(testUserName)
	if err != nil {
		t.Fatalf("error getting user: %s\n", err)
	}

	if user.FirstName != testFirstName {
		t.Fatalf("first name does not match: %s and %s\n", user.FirstName, testFirstName)
	}

	if user.LastName != testLastName {
		t.Fatalf("last name does not match: %s and %s\n", user.LastName, testLastName)
	}

	_, err = service.User(ewserver.UserName("notexist"))
	if err == nil {
		t.Fatalf("error should have got error when user does not exist")
	}

	if err != ewserver.ErrUserNotFound {
		t.Fatalf("error should have been not found, got: %s\n", err)
	}
}

func TestUserService_Users(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	service := boltdb.NewUserService(db.DB())
	if err := service.Init(); err != nil {
		t.Fatalf("error initializing user service: %s\n", err)
	}

	users, err := service.Users()
	if err != nil {
		t.Fatalf("got error when users is empty: %s\n", err)
	}

	if len(users) != 0 {
		t.Fatalf("users should be empty")
	}

	testCreateUsers(service, t)

	users, err = service.Users()
	if err != nil {
		t.Fatalf("error getting users: %s\n", err)
	}

	if len(users) != 10 {
		t.Fatalf("expected 10 users got: %d\n", len(users))
	}
}

func TestUserService_Authenticate(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	service := boltdb.NewUserService(db.DB())
	if err := service.Init(); err != nil {
		t.Fatalf("error initializing user service: %s\n", err)
	}

	testCreateUser(service, t)

	user, err := service.User(ewserver.UserName(testUserName))
	if err != nil {
		t.Fatalf("error getting users: %s\n", err)
	}

	// try correct password
	authUser, err := service.Authenticate(user.UserName, testPassword)
	if err != nil {
		t.Fatalf("error authenticating user: %s\n", err)
	}

	if authUser.FirstName != user.FirstName {
		t.Fatalf("first names don't match got %s expected %s", testFirstName, authUser.FirstName)
	}

	// try wrong
	authUser, err = service.Authenticate(user.UserName, "wrong")
	if err != ewserver.ErrInvalidPassword {
		t.Fatalf("error should have got invalid password error, got: %s\n", err)
	}

	// try empty
	authUser, err = service.Authenticate(user.UserName, "")
	if err != ewserver.ErrInvalidPassword {
		t.Fatalf("error should have got invalid password error, got: %s\n", err)
	}

	// try invalid username
	authUser, err = service.Authenticate("", "")
	if err != ewserver.ErrUserNotFound {
		t.Fatalf("error should have got user not found error, got: %s\n", err)
	}

	// try different real users password
	u := ewserver.NewUser()
	u.UserName = "different"
	if err := service.Create(u, "blonk"); err != nil {
		t.Fatalf("error creating new user: %s\n", err)
	}

	authUser, err = service.Authenticate(u.UserName, testPassword)
	if err != ewserver.ErrInvalidPassword {
		t.Fatalf("error should have got invalid password error, got: %s\n", err)
	}
}

func TestUserService_ChangePassword(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	service := boltdb.NewUserService(db.DB())
	if err := service.Init(); err != nil {
		t.Fatalf("error initializing user service: %s\n", err)
	}

	// try not in database yet
	if err := service.ChangePassword(testUserName, testPassword, "not in db yet"); err != nil && err != ewserver.ErrUserNotFound {
		t.Fatalf("error changing password expected not found got: %s\n", err)
	}

	testCreateUser(service, t)

	user, err := service.User(ewserver.UserName(testUserName))
	if err != nil {
		t.Fatalf("error getting user: %s\n", err)
	}

	if err := service.ChangePassword(user.UserName, testPassword, "newpassword"); err != nil {
		t.Fatalf("error changing password: %s\n", err)
	}

	// try old password
	if _, err := service.Authenticate(user.UserName, testPassword); err != ewserver.ErrInvalidPassword {
		t.Fatalf("error old password should have returned invalid password error got: %s\n", err)
	}

	// try new password
	if _, err := service.Authenticate(user.UserName, "newpassword"); err != nil {
		t.Fatalf("error new password did not authenticate user: %s\n", err)
	}

	// try changing password with incorrect/empty password
	if err := service.ChangePassword(user.UserName, "", "newpassword"); err != ewserver.ErrInvalidPassword {
		t.Fatalf("error should have returned invalid password error got: %s\n", err)
	}
}

func TestUserService_ResetPassword(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	service := boltdb.NewUserService(db.DB())
	if err := service.Init(); err != nil {
		t.Fatalf("error initializing user service: %s\n", err)
	}

	// try not in database yet
	if err := service.ResetPassword(testUserName, "not in db yet"); err != nil && err != ewserver.ErrUserNotFound {
		t.Fatalf("error changing password expected not found got: %s\n", err)
	}

	testCreateUser(service, t)

	user, err := service.User(ewserver.UserName(testUserName))
	if err != nil {
		t.Fatalf("error getting user: %s\n", err)
	}

	if err := service.ResetPassword(user.UserName, "newpassword"); err != nil {
		t.Fatalf("error changing password: %s\n", err)
	}

	// try old password
	if _, err := service.Authenticate(user.UserName, testPassword); err != ewserver.ErrInvalidPassword {
		t.Fatalf("error old password should have returned invalid password error got: %s\n", err)
	}

	// try new password
	if _, err := service.Authenticate(user.UserName, "newpassword"); err != nil {
		t.Fatalf("error new password did not authenticate user: %s\n", err)
	}
}

func TestUserService_Delete(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	service := boltdb.NewUserService(db.DB())
	if err := service.Init(); err != nil {
		t.Fatalf("error initializing user service: %s\n", err)
	}

	if err := service.Delete(testUserName); err != nil {
		t.Fatalf("error should not have got any errors when deleting non-existent user: %s\n", err)
	}

	testCreateUser(service, t)

	_, err = service.User(testUserName)
	if err != nil {
		t.Fatalf("error getting existing user: %s\n", err)
	}

	if err := service.Delete(testUserName); err != nil {
		t.Fatalf("error should not have got any errors when deleting an existing user: %s\n", err)
	}

	if _, err = service.User(testUserName); err != nil && err != ewserver.ErrUserNotFound {
		t.Fatalf("error did not get user not found error, got: %s\n", err)
	}

}

func testCreateUser(service *boltdb.UserService, t *testing.T) {
	u := ewserver.NewUser()
	u.UserName = testUserName
	u.FirstName = testFirstName
	u.LastName = testLastName

	if err := service.Create(u, testPassword); err != nil {
		t.Fatalf("error creating user: %s\n", err)
	}
}

func testCreateUsers(service *boltdb.UserService, t *testing.T) {
	for i := 0; i < 10; i++ {
		u := ewserver.NewUser()
		u.UserName = ewserver.UserName(fmt.Sprintf("%s%d", testUserName, i))
		u.FirstName = testFirstName
		u.LastName = testLastName

		if err := service.Create(u, testPassword); err != nil {
			t.Fatalf("error creating user: %s\n", err)
		}
	}
}
