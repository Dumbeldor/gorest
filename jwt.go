package gorest

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Session struct {
	Username string
	UserID   int64
}

func CreateJWT(username string, userID int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	return token.SignedString(GetSecret())
}

func GetSession(c echo.Context) *Session {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return &Session{
		claims["username"].(string),
		int64(claims["id"].(float64)),
	}
}

func GetSecret() []byte {
	return []byte(fmt.Sprint(Get("secret")))
}
