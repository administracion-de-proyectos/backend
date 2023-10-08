package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

type TokenValidator[T any] interface {
	CreateToken(data T) (string, error)
	ValidateToken(token string) error
	GetData(token string) (T, error)
}

type TokenValidatorRequest interface {
	ValidateTokenMiddleware(c *gin.Context)
}
type tokenValidator[T any] struct {
	secret string
}

func (t *tokenValidator[T]) CreateToken(data T) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"data": data,
	})
	log.Infof("secret is: %s", t.secret)
	s, err := token.SignedString([]byte(t.secret))
	if err != nil {
		log.Errorf("some fuckery happened with the jwt: %s", err.Error())
	}
	return s, err
}

func (t *tokenValidator[T]) ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(t.secret), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

type MyCustomClaims[T any] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}

func (t *tokenValidator[T]) GetData(tokenString string) (T, error) {
	var data MyCustomClaims[T]
	_, err := jwt.ParseWithClaims(tokenString, &data, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(t.secret), nil
	})
	if err != nil {
		return data.Data, err
	}
	res := data.Data
	return res, err
}

func CreateValidator[T any](secret string) TokenValidator[T] {
	return &tokenValidator[T]{secret: secret}
}
