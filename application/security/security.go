package security

import (
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"time"

	"jsotogaviard-api-test/application/constants"
	"errors"
)

type Security struct{
	Jwt *jwtmiddleware.JWTMiddleware
}

// TODO Read key from file
var signingKey = []byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9")

// Init security
func GetSecurity() *Security{
	var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return &Security{jwtMiddleware}
}

// Compute the token
func ComputeToken(login string, id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims[constants.GetLogin()] = login
	claims[constants.GetUserId()] = id
	claims[constants.GetExpiration()] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

func GetUserId(tokenString string) (*float64, error) {

	// Check signing method
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method: " +  token.Header["alg"].(string))
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token has not expired
		now := time.Now().Unix()
		expiration := claims[constants.GetExpiration()]
		expirationTime := int64(expiration.(float64))
		if now > expirationTime {
			return nil, errors.New("Expired token")
		} else {
			userId := claims[constants.GetUserId()]
			userIdInt := userId.(float64)
			return &userIdInt, nil
		}
	} else {
		return nil, errors.New("Not valid token")

	}
}
