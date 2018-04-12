package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"github.com/wirepair/ewserver/store/boltdb"

	"github.com/gin-gonic/gin"

	"github.com/wirepair/ewserver/api/v1"
	"golang.org/x/crypto/acme/autocert"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config/server.json", "path to server config json file")
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

	gin.SetMode(gin.DebugMode)
	e := gin.Default()
	routes := e.Group("v1")
	routes.Use()
	v1.RegisterAdminRoutes(userService, routes, e)
	v1.RegisterAdminAPIRoutes(apiUserService, routes, e)

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
