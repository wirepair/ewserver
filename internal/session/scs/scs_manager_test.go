package scs

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/wirepair/ewserver/ewserver"
	"github.com/wirepair/ewserver/store"
	"github.com/wirepair/ewserver/store/boltdb"
	"github.com/wirepair/scs/stores/boltstore"

	"github.com/alexedwards/scs"
)

func TestSessions_AddString(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)
	data := db.DB()
	s := boltstore.New(data, time.Hour*12)

	manager := scs.NewManager(s)
	sessions := New(manager)

	// build add string request
	req := httptest.NewRequest("GET", "http://ewserver", nil)
	w := httptest.NewRecorder()
	add := testAddStringHandler(sessions, t)
	add.ServeHTTP(w, req)

	cookie := testExtractCookie(w, t)

	// send next request with cookie
	req = httptest.NewRequest("GET", "http://ewserver", nil)
	req.Header.Add("cookie", cookie)
	w = httptest.NewRecorder()

	// get and test string is returned
	get := testGetStringHandler(sessions, t)
	get.ServeHTTP(w, req)
}

func TestSessions_AddObject(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)
	data := db.DB()
	s := boltstore.New(data, time.Hour*12)

	manager := scs.NewManager(s)
	sessions := New(manager)

	req := httptest.NewRequest("GET", "http://ewserver", nil)
	w := httptest.NewRecorder()
	add := testAddObjectHandler(sessions, t)
	add.ServeHTTP(w, req)

	cookie := testExtractCookie(w, t)

	req = httptest.NewRequest("GET", "http://ewserver", nil)
	req.Header.Add("cookie", cookie)
	w = httptest.NewRecorder()

	get := testGetObjectHandler(sessions, t)
	get.ServeHTTP(w, req)
}

func TestSessions_Renew(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)
	data := db.DB()
	s := boltstore.New(data, time.Hour*12)

	manager := scs.NewManager(s)
	sessions := New(manager)

	req := httptest.NewRequest("GET", "http://ewserver", nil)
	w := httptest.NewRecorder()
	add := testAddStringHandler(sessions, t)
	add.ServeHTTP(w, req)

	cookie := testExtractCookie(w, t)
	t.Logf("Value: %s\n", cookie)

	req = httptest.NewRequest("GET", "http://ewserver", nil)
	req.Header.Add("cookie", cookie)
	w = httptest.NewRecorder()

	renew := testRenewHandler(sessions, t)
	renew.ServeHTTP(w, req)

	cookie2 := testExtractCookie(w, t)
	t.Logf("Value: %s\n", cookie2)
}

func TestSessions_PopString(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)
	data := db.DB()
	s := boltstore.New(data, time.Hour*12)

	manager := scs.NewManager(s)
	sessions := New(manager)

	req := httptest.NewRequest("GET", "http://ewserver", nil)
	w := httptest.NewRecorder()
	add := testAddStringHandler(sessions, t)
	add.ServeHTTP(w, req)

	cookie := testExtractCookie(w, t)

	req = httptest.NewRequest("GET", "http://ewserver", nil)
	req.Header.Add("cookie", cookie)
	w = httptest.NewRecorder()

	get := testPopStringHandler(sessions, t)
	get.ServeHTTP(w, req)

	cookie = testExtractCookie(w, t)

	req = httptest.NewRequest("GET", "http://ewserver", nil)
	req.Header.Add("cookie", cookie)
	w = httptest.NewRecorder()

	get = testGetStringEmptyHandler(sessions, t)
	get.ServeHTTP(w, req)
}

func TestSessions_PopObject(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)
	data := db.DB()
	s := boltstore.New(data, time.Hour*12)

	manager := scs.NewManager(s)
	sessions := New(manager)

	req := httptest.NewRequest("GET", "http://ewserver", nil)
	w := httptest.NewRecorder()
	add := testAddObjectHandler(sessions, t)
	add.ServeHTTP(w, req)

	cookie := testExtractCookie(w, t)

	req = httptest.NewRequest("GET", "http://ewserver", nil)
	req.Header.Add("cookie", cookie)
	w = httptest.NewRecorder()

	get := testPopObjectHandler(sessions, t)
	get.ServeHTTP(w, req)

	cookie = testExtractCookie(w, t)

	req = httptest.NewRequest("GET", "http://ewserver", nil)
	req.Header.Add("cookie", cookie)
	w = httptest.NewRecorder()

	get = testGetObjectEmptyHandler(sessions, t)
	get.ServeHTTP(w, req)
}

func testExtractCookie(w *httptest.ResponseRecorder, t *testing.T) string {
	cookies := w.Header().Get("set-cookie")
	if cookies == "" {
		t.Fatalf("cookie was not set\n")
	}

	cookie := strings.Split(cookies, " ")
	return cookie[0]
}

func testRenewHandler(sessions *Sessions, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		sessions.Renew(w, req)
		io.WriteString(w, "test")
	}
}

func testAddObjectHandler(sessions *Sessions, t *testing.T) http.HandlerFunc {
	apiUser := &ewserver.APIUser{}
	apiUser.ID = []byte("123")
	//gob.Register(&apiUser)
	return func(w http.ResponseWriter, req *http.Request) {
		sessions.Add(w, req, "test", apiUser)
		io.WriteString(w, "test")
	}
}

func testGetObjectHandler(sessions *Sessions, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		apiUser := &ewserver.APIUser{}
		err := sessions.Load(req, "test", apiUser)
		if err != nil {
			t.Fatalf("error loading object: %s\n", err)
		}
		if bytes.Compare(apiUser.ID, []byte("123")) != 0 {
			t.Fatalf("expected 123 got: %#v\n", apiUser.ID)
		}
		io.WriteString(w, "test")
	}
}

func testAddStringHandler(sessions *Sessions, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		sessions.Add(w, req, "test", "blah")
		io.WriteString(w, "test")
	}
}

func testGetStringHandler(sessions *Sessions, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result := sessions.GetString(w, req, "test")
		if result != "blah" {
			t.Fatalf("expected %s got %s\n", "blah", result)
		}
		io.WriteString(w, "test")
	}
}

func testGetStringEmptyHandler(sessions *Sessions, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result := sessions.GetString(w, req, "test")
		if result != "" {
			t.Fatalf("expected empty string got %s\n", result)
		}
		io.WriteString(w, "test")
	}
}

func testGetObjectEmptyHandler(sessions *Sessions, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		apiUser := &ewserver.APIUser{}
		err := sessions.Load(req, "test", apiUser)
		if err != nil {
			t.Fatalf("error loading object: %s\n", err)
		}

		if apiUser.ID != nil {
			t.Fatalf("expected id to be nil, got %v\n", apiUser.ID)
		}
		io.WriteString(w, "test")
	}
}

func testPopStringHandler(sessions *Sessions, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result := sessions.PopString(w, req, "test")
		if result != "blah" {
			t.Fatalf("expected %s got %s\n", "blah", result)
		}
		io.WriteString(w, "test")
	}
}

func testPopObjectHandler(sessions *Sessions, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		apiUser := &ewserver.APIUser{}
		err := sessions.PopLoad(w, req, "test", apiUser)
		if err != nil {
			t.Fatalf("error loading object: %s\n", err)
		}
		if bytes.Compare(apiUser.ID, []byte("123")) != 0 {
			t.Fatalf("expected 123 got: %#v\n", apiUser.ID)
		}
		io.WriteString(w, "test")
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
