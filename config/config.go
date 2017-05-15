package config

import (
	"github.com/spankie/web-chat/models"
	// "github.com/tidwall/buntdb"
)

// Config holds the application wide data.
type Config struct {
	DB map[int]models.User
}

var (
	config Config
)

// init initialises the config
func init() {
	config = Config{DB: make(map[int]models.User, 10)}

}

// Get returns the config
func Get() *Config {
	return &config
}
