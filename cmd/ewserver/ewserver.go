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
	"github.com/gin-gonic/gin"
	"github.com/wirepair/ewserver/api/v1"
	"github.com/wirepair/ewserver/api/v1/middleware"
	"github.com/wirepair/ewserver/internal/authz/casbinauth"
	"github.com/wirepair/ewserver/internal/session/scssession"
	"github.com/wirepair/ewserver/store/boltdb"
	"github.com/wirepair/scs/stores/boltstore"
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
	manager := scs.NewManager(sessionStore)
	sessions := scssession.New(manager)

	// initialize authz
	boltauth := boltadapter.NewAdapter(db.DB())
	enforcer := casbin.NewEnforcer(serverConfig.AuthPolicyPath, boltauth)
	authorizer := casbinauth.New(enforcer, apiUserService, sessions)

	if debug {
		gin.SetMode(gin.DebugMode)
		// allow admin access to everything under admin
		enforcer.AddPolicy("admin", "/v1/api/admin/", "*")
		// add the testuser to the apiusers role
		enforcer.AddGroupingPolicy("admin", "admin")
		boltauth.SavePolicy(enforcer.GetModel())
	}

	// setup server
	e := gin.Default()
	e.Use(middleware.EnsureSession(sessions))

	routes := e.Group("v1")
	v1.RegisterAuthnRoutes(userService, routes, e)

	adminRoutes := routes.Group("/admin/users")
	adminRoutes.Use(middleware.Require(authorizer))
	v1.RegisterAdminRoutes(userService, adminRoutes, e)

	apiAdminRoutes := routes.Group("/admin/api_users")
	v1.RegisterAdminAPIRoutes(apiUserService, apiAdminRoutes, e)

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
