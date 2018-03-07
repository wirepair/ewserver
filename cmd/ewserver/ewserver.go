package server

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/wirepair/ewserver/api/v1"
	"github.com/wirepair/ewserver/config"
	"github.com/wirepair/ewserver/store"
	"golang.org/x/crypto/acme/autocert"
)

// main runs the server
func main(config *config.ServerConfig) error {
	store, err := store.NewStore(config.StoreEngine, config.StoreConfig)
	if err != nil {
		return err
	}
	gin.SetMode(gin.DebugMode)
	e := gin.Default()

	v1.RegisterRoutes(store, e)

	if config.EnableHTTPS {
		createCacheDir(config)
		go log.Fatal(runWithManager(e, config))
	}

	log.Fatal(e.Run(config.HTTPAddr))
	return nil
}

// runWithManager starts an https server with lets encrypt / acme support.
// Note TLS port *must* be 443 if lets encrypt.
func runWithManager(e *gin.Engine, serverConfig *config.ServerConfig) error {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(serverConfig.Domains...),
		Cache:      autocert.DirCache(serverConfig.CacheDir),
	}

	s := &http.Server{
		Addr:      serverConfig.HTTPSAddr,
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   e,
	}
	return s.ListenAndServeTLS("", "")
}

// createCacheDir for lets encrypt
func createCacheDir(config *config.ServerConfig) {
	if fileInfo, err := os.Stat(config.CacheDir); err == nil {
		if fileInfo.IsDir() {
			return
		}
		log.Fatalf("error cache folder is not a directory: %s\n", err)
	}

	if err := os.MkdirAll(config.CacheDir, 0700); err != nil {
		log.Fatalf("error creating cache directory: %s\n", err)
	}
}
