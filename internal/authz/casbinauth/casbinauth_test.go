package casbinauth

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/wirepair/bolt-adapter"

	"github.com/casbin/casbin"
	"github.com/wirepair/ewserver/ewserver"
	"github.com/wirepair/ewserver/mock"
	"github.com/wirepair/ewserver/store"
	"github.com/wirepair/ewserver/store/boltdb"
)

var logger = &mock.Log{}

func TestCasbinAuthorizer_Authorize(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	adapter := boltadapter.NewAdapter(db.DB())
	enforcer := casbin.NewSyncedEnforcer("testdata/rbac_model.conf", adapter)

	// allow 'apiusers' role to access /api/ via GET
	enforcer.AddPolicy("apiusers", "/api/", "GET")

	// add the testuser to the apiusers role
	enforcer.AddGroupingPolicy("testuser", "apiusers")

	adapter.SavePolicy(enforcer.GetModel())

	sessions := &mock.Sessions{}
	sessions.LoadFn = func(req *http.Request, key string, val interface{}) error {
		if user, ok := val.(*ewserver.User); ok {
			user.UserName = "testuser"
		}
		return nil
	}
	user := &ewserver.User{}

	usapi := &mock.APIUserService{}
	auth := NewAuthorizer(enforcer, usapi, sessions, logger)
	req := httptest.NewRequest("GET", "http://ewserver/api/", nil)

	sessions.LoadFn(req, "test", user)
	if user.UserName != "testuser" {
		t.Fatalf("error not testuser\n")
	}

	if !auth.Authorize(req) {
		t.Fatalf("error GET should be authorized\n")
	}

	req = httptest.NewRequest("GET", "http://ewserver/api/asdf", nil)
	if !auth.Authorize(req) {
		t.Fatalf("error GET /api/* should be authorized\n")
	}

	req = httptest.NewRequest("POST", "http://ewserver/api/", nil)
	if auth.Authorize(req) {
		t.Fatalf("error POST /api/ be denied\n")
	}

	req = httptest.NewRequest("GET", "http://ewserver/notapi", nil)
	if auth.Authorize(req) {
		t.Fatalf("error GET /notapi/ should be denied\n")
	}

	/*
		// This is scary that it passes, must confirm it does NOT authorize in a real request.
		req = httptest.NewRequest("GET", "http://ewserver/api/../notapi", nil)
		if auth.Authorize(req) {
			t.Fatalf("error GET /notapi/ (via traversal) should be denied.\n")
		}
	*/
}

func TestCasbinAuthorizer_Authorize2(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	adapter := boltadapter.NewAdapter(db.DB())
	enforcer := casbin.NewSyncedEnforcer("testdata/rbac_model.conf", adapter)
	enforcer.AddPolicy("admin", "/", ".*")
	enforcer.AddPolicy("apiuser", "/v1/api/", "(GET|POST)")
	// only allow anonymous to access the top folder
	enforcer.AddPolicy("anonymous", "/:", "(GET|POST)")
	// add root to the admin role
	enforcer.AddGroupingPolicy("root", "admin")
	sessions := &mock.Sessions{}
	sessions.LoadFn = func(req *http.Request, key string, val interface{}) error {
		if user, ok := val.(*ewserver.User); ok {
			user.UserName = "anonymous"
		}
		return nil
	}
	user := &ewserver.User{}
	usapi := &mock.APIUserService{}
	auth := NewAuthorizer(enforcer, usapi, sessions, logger)
	req := httptest.NewRequest("GET", "http://ewserver/v1/admin/users/all_details", nil)

	sessions.LoadFn(req, "test", user)

	if auth.Authorize(req) {
		t.Fatalf("error GET /v1/admin/users/all_details should not be authorized\n")
	}

	req = httptest.NewRequest("GET", "http://ewserver/static/", nil)
	sessions.LoadFn(req, "test", user)
	if auth.Authorize(req) {
		t.Fatalf("error GET /static/ should not be authorized\n")
	}

	req = httptest.NewRequest("GET", "http://ewserver/login", nil)
	sessions.LoadFn(req, "test", user)
	if !auth.Authorize(req) {
		t.Fatalf("error GET /login SHOULD be authorized\n")
	}
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
