package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/wirepair/ewserver/store"
)

const serverFile = "server.json"

// ServerConfig holds various configuration data for the top level service
type ServerConfig struct {
	CacheDir       string        `json:"cache_dir"`       // for lets encrypt
	StoreEngine    string        `json:"store_engine"`    // the store engine to use: ql, postgres, mysql
	StoreConfig    *store.Config `json:"store_config"`    // the store configuration options
	Host           string        `json:"host"`            // Our hostname or IP address
	SMTPServer     string        `json:"smtp_server"`     // host:port of an SMTP server to send notification emails through
	UseLetsEncrypt bool          `json:"use_letsencrypt"` // use lets encrypt based TLS.
	EnableHTTPS    bool          `json:"enable_https"`    // if we want to enable https + letsencrypt
	HTTPAddr       string        `json:"http_addr"`       // the http address to bind to, like :8080
	HTTPSAddr      string        `json:"https_addr"`      // the https address to bind to, like :8443
}

// ReadServerConfig reads the server config from a json file.
func ReadServerConfig(path string) *ServerConfig {
	file, err := os.OpenFile(path+string(os.PathSeparator)+serverFile, os.O_RDONLY, 0600)
	if err != nil {
		log.Fatalf("error getting server config: %s\n", err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("error reading server file: %s\n", err)
	}

	serverConfig := &ServerConfig{}
	if err := json.Unmarshal(data, serverConfig); err != nil {
		log.Fatalf("error unmarshalling json server config: %s\n", err)
	}

	return serverConfig
}
