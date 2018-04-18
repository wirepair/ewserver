package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/wirepair/bolt-adapter"

	"github.com/casbin/casbin"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/boltstore"
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/api/v1"
	"github.com/wirepair/ewserver/api/v1/middleware"
	"github.com/wirepair/ewserver/ewserver"
	"github.com/wirepair/ewserver/internal/authz/casbinauth"
	"github.com/wirepair/ewserver/internal/session/scssession"
	"github.com/wirepair/ewserver/store/boltdb"
	"golang.org/x/crypto/acme/autocert"
)

var configPath string
var debug bool

func init() {
	flag.StringVar(&configPath, "config", "config/server.json", "path to server config json file")
	flag.BoolVar(&debug, "debug", true, "debug mode")
}

// main runs the HTTP(s) server
func main() {
	flag.Parse()

	serverConfig := ReadServerConfig(configPath)

	db := boltdb.NewBoltStore()
	if err := db.Open(serverConfig.StoreConfig); err != nil {
		log.Fatalf("error opening config: %s\n", err)
	}

	userService := boltdb.NewUserService(db.DB())
	if err := userService.Init(); err != nil {
		log.Fatalf("error initializing UserService: %s\n", err)
	}

	apiUserService := boltdb.NewAPIUserService(db.DB())
	if err := apiUserService.Init(); err != nil {
		log.Fatalf("error initializing APIUserService: %s\n", err)
	}

	// initialize sessions
	sessionStore := boltstore.New(db.DB(), time.Minute)
	scsManager := scs.NewManager(sessionStore)
	sessions := scssession.New(scsManager)

	// initialize authz
	boltauth := boltadapter.NewAdapter(db.DB())
	enforcer := casbin.NewSyncedEnforcer(serverConfig.AuthPolicyPath, boltauth)
	authorizer := casbinauth.NewAuthorizer(enforcer, apiUserService, sessions)

	roleService := casbinauth.NewRoleService(enforcer)
	services := ewserver.NewServices(userService, apiUserService, roleService)

	if debug {
		gin.SetMode(gin.DebugMode)
		// allow admin access to everything
		enforcer.AddPolicy("admin", "/", ".*")
		enforcer.AddPolicy("apiuser", "/v1/api/:", ".*")
		// only allow anonymous to access the top folder
		enforcer.AddPolicy("anonymous", "/:", "(GET|POST)")
		// add root to the admin role
		enforcer.AddGroupingPolicy("root", "admin")
		boltauth.SavePolicy(enforcer.GetModel())
		root := &ewserver.User{UserName: "root"}
		userService.Create(root, "password")
	}

	// setup server
	e := gin.Default()
	e.Use(middleware.EnsureSession(sessions), middleware.Require(authorizer))

	v1.RegisterAuthnRoutes(userService, e)
	v1.RegisterAdminRoutes(services, e)

	if serverConfig.EnableHTTPS {
		go log.Fatal(runWithManager(e, serverConfig))
	}

	log.Fatal(e.Run(serverConfig.HTTPAddr))
}

// runWithManager starts an https server with lets encrypt / acme support.
// Note TLS port *must* be 443 if lets encrypt.
func runWithManager(e *gin.Engine, serverConfig *ServerConfig) error {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(serverConfig.Host),
		Cache:      autocert.DirCache(serverConfig.CacheDir),
	}

	s := &http.Server{
		Addr:      serverConfig.HTTPSAddr,
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   e,
	}
	return s.ListenAndServeTLS("", "")
}
