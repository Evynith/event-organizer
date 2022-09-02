package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

//using HMAC!!!

// jwt service
type JWTService interface {
	GenerateToken(email string) string
	ValidateToken(token string) (*jwt.Token, error)
	TypeUser(encodedToken string) (string, string)
}
type authCustomClaims struct {
	Name string `json:"name"`
	User string `json:"user"`
	Id   string `json:"id"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issure    string
}

// auth-jwt
func JWTAuthService() JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issure:    "Rohan",
	}
}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func (service *jwtServices) GenerateToken(username string) string {
	user, id := DataUser(username)
	claims := &authCustomClaims{
		username,
		user,
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    service.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])

		}
		return []byte(service.secretKey), nil
	})

}

func (service *jwtServices) TypeUser(encodedToken string) (string, string) {
	token, err := service.ValidateToken(encodedToken)
	if err != nil {
		return "", ""
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims["user"].(string), claims["id"].(string)
		}
	}
	return "", ""
}
