package web

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spankie/web-chat/config"
	"github.com/spankie/web-chat/models"
)

// GenerateJWT Turns user details into a hashed token that can be used to recognize the user in the future.
func GenerateJWT(user models.User) (tokenString string, err error) {

	claims := jwt.MapClaims{}

	// set our claims
	claims["User"] = user
	claims["Name"] = user.Username

	// set the expire time
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30 * 12).Unix() //24 hours in a day, in 30 days * 12 months = 1 year in milliseconds

	// create a signer for rsa 256
	t := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)

	pub, err := jwt.ParseRSAPrivateKeyFromPEM(config.Get().PrivateKey)
	if err != nil {
		return
	}

	tokenString, err = t.SignedString(pub)

	if err != nil {
		return
	}

	return

}

func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
