package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

//using HMAC!!!

type JWTService interface {
	GenerateToken(username string) string
	ValidateToken(token string) (*jwt.Token, error)
	TypeUser(parseToken *jwt.Token) (string, string)
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

/*
Recibe un nombre de usuario que usará para buscar el resto de datos y generar un token a base de dichos datos
por un tiempo limitado de 48hs
*/
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

	t, err := token.SignedString([]byte(service.secretKey)) //encoded string
	if err != nil {
		panic(err)
	}
	return t
}

/*
Revisa que el token recibido tenga la estructura correcta a través de su decodificación y sea válido según sus claims de creación,
luego verifica que se  corresponde  con la firma de creación guardada en el servicio. Devuelve el token parseado según dichas claims
*/
func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])

		}
		return []byte(service.secretKey), nil
	})

}

/*
Recibe un token validado y parseado del que retorna el contenido asociado a las claims de "user" y de "id"
*/
func (service *jwtServices) TypeUser(parseToken *jwt.Token) (string, string) {
	if claims, ok := parseToken.Claims.(jwt.MapClaims); ok {
		return claims["user"].(string), claims["id"].(string)
	}
	return "", ""
}
