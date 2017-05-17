package config

import (
	"io/ioutil"
	"log"

	"github.com/spankie/web-chat/models"
	// "github.com/tidwall/buntdb"
)

// Config holds the application wide data.
type Config struct {
	DB         map[int]models.User
	PrivateKey []byte
}

var (
	config Config
)

// init initialises the config
func init() {
	config = Config{DB: make(map[int]models.User, 10)}

	// Just to make PrivateKey assign on the next line
	var err error

	config.PrivateKey, err = ioutil.ReadFile("./config/keys/key.pem")
	if err != nil {
		log.Println("Error reading public key")
		log.Println(err)
		return
	}

}

// Get returns the config
func Get() *Config {
	return &config
}
