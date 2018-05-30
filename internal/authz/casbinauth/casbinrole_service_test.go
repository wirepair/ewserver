package casbinauth

import (
	"reflect"
	"testing"

	"github.com/casbin/casbin"
	boltadapter "github.com/wirepair/bolt-adapter"
)

func TestCasbinRoleService_Users(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	adapter := boltadapter.NewAdapter(db.DB())
	enforcer := casbin.NewSyncedEnforcer("testdata/rbac_model.conf", adapter)
	testAddDefaultPolicy(enforcer)
	service := NewRoleService(enforcer)

	roleNames := service.RoleNames()
	if len(roleNames) != 3 {
		t.Fatalf("expected 3 users got: %d\n", len(roleNames))
	}
}

func TestCasbinRoleService_Groups(t *testing.T) {
	dbFileName, err := testTempDbFileName("testdata/")
	if err != nil {
		t.Fatalf("error opening db file for testing")
	}
	defer testRemoveDbFile(dbFileName, t)

	db := testOpenDb(dbFileName, t)
	defer testCloseDb(db, t)

	adapter := boltadapter.NewAdapter(db.DB())
	enforcer := casbin.NewSyncedEnforcer("testdata/rbac_model.conf", adapter)
	testAddDefaultPolicy(enforcer)
	service := NewRoleService(enforcer)

	roleMap := service.RoleMap()
	if len(roleMap) != 1 {
		t.Fatalf("expected 1 role got: %d\n", len(roleMap))
	}

	if !reflect.DeepEqual(roleMap[0], []string{"root", "admin"}) {
		t.Fatalf("expected roleMap to contain admin -> root mapping: got %#v\n", roleMap[0])
	}

	if err := service.DeleteRole("admin"); err != nil {
		t.Fatalf("unable to delete group: ")
	}
}

func testAddDefaultPolicy(enforcer *casbin.SyncedEnforcer) {
	enforcer.AddPolicy("admin", "/", ".*")
	enforcer.AddPolicy("apiuser", "/v1/api/", "(GET|POST)")
	// only allow anonymous to access the top folder
	enforcer.AddPolicy("anonymous", "/:", "(GET|POST)")
	// add root to the admin role
	enforcer.AddGroupingPolicy("root", "admin")
}
