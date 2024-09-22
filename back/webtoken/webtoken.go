package webtoken

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"

	"github.com/code-troopers/postitsonline/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func Middleware(pubKey *rsa.PublicKey) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Token manquant")
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		_, err := ValidateToken(token, pubKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Token invalide")
		}
		user, err := DecodeJWT(token)
		if err != nil {
			log.Errorf("Erreur lors du décodage du token : %v", err)
		}
		c.Locals("user", user)

		return c.Next()
	}
}

func ValidateToken(tokenString string, pubKey *rsa.PublicKey) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil || !token.Valid {
		log.Errorf("Token invalide : %v", err)
		return nil, err
	}
	return token, nil
}

func DecodeJWT(tokenString string) (database.User, error) {
	publicKey, err := GetKeycloakPublicKey()
	if err != nil {
		return database.User{}, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("méthode de signature invalide : %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return database.User{}, fmt.Errorf("erreur lors du décodage du token : %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["sub"].(string)
		if !ok {
			return database.User{}, fmt.Errorf("sub non trouvé dans le token")
		}
		givenName, ok := claims["given_name"].(string)
		if !ok {
			return database.User{}, fmt.Errorf("givenName non trouvé dans le token")
		}
		familyName, ok := claims["family_name"].(string)
		if !ok {
			return database.User{}, fmt.Errorf("familyName non trouvé dans le token")
		}
		email, ok := claims["email"].(string)
		if !ok {
			return database.User{}, fmt.Errorf("email non trouvé dans le token")
		}
		picture, ok := claims["picture"].(string)
		if !ok {
			return database.User{}, fmt.Errorf("picture non trouvé dans le token")
		}

		user := database.User{
			ID:         id,
			GivenName:  givenName,
			FamilyName: familyName,
			Email:      email,
			Picture:    picture,
		}
		return user, nil
	}

	return database.User{}, fmt.Errorf("token invalide ou signature incorrecte")
}
func GetKeycloakPublicKey() (*rsa.PublicKey, error) {
	keycloakCerts := os.Getenv("KEYCLOAK_CERTS")
	if keycloakCerts == "" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
		keycloakCerts = os.Getenv("KEYCLOAK_CERTS")
	}
	resp, err := http.Get(keycloakCerts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwks struct {
		Keys []struct {
			N string `json:"n"` // Modulus
			E string `json:"e"` // Exponent
		} `json:"keys"`
	}

	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, err
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(jwks.Keys[0].N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(jwks.Keys[0].E)
	if err != nil {
		return nil, err
	}

	e := big.NewInt(0)
	e.SetBytes(eBytes)

	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(e.Int64()),
	}

	return pubKey, nil
}
