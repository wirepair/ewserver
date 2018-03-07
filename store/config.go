package store

// Config is a generic config file
type Config struct {
	Options map[string]string `json:"options"`
}

// NewConfig for an arbitrary data store
func NewConfig() *Config {
	c := &Config{}
	c.Options = make(map[string]string, 0)
	return c
}
