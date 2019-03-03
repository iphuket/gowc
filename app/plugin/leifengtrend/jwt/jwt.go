package jwt

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/iphuket/gowc/app/plugin/leifengtrend/crypto"
)

// ParsingToken jwt token string
func ParsingToken() (string, error) {

	return "", nil
}

// Aeskey ... JWT
var Aeskey = "ekwjunfw87&&*12"

// NewToken new jwt token string
func NewToken(uuid, sub, ip string, sec int64) (string, error) {
	// encrypt user info
	enuuid, err := crypto.EnSting([]byte(uuid), []byte(Aeskey))
	if err != nil {
		return "aes error ", err
	}
	ensub, err := crypto.EnSting([]byte(sub), []byte(Aeskey))
	if err != nil {
		return "aes error ", err
	}
	enip, err := crypto.EnSting([]byte(ip), []byte(Aeskey))
	if err != nil {
		return "aes error ", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": enuuid,
		"sub":  ensub,
		"ip":   enip,
		"exp":  time.Now().Unix() + sec,
		"nbf":  time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(Aeskey))
	if err != nil {
		return "Signing error ", err
	}
	return tokenString, nil
}

// CheckIP ip check
func CheckIP(ip string, Token *jwt.Token) (bool, error) {
	if claims, ok := Token.Claims.(jwt.MapClaims); ok {
		ip, err := crypto.EnSting([]byte(ip), []byte(Aeskey))
		if err != nil {
			return false, err
		}
		if claims["ip"].(string) != string(ip) {
			return false, err
		}
		return true, nil
	}
	return false, errors.New("claims false")
}

// Chcek chcek user login status and Get the UUID. auto chcek ip
func Chcek(ip string, token string) (string, error) {
	Token, err := jwt.Parse(token, parseInterface)
	if err != nil {
		return "", err
	}
	if claims, ok := Token.Claims.(jwt.MapClaims); ok && Token.Valid {
		bool, err := CheckIP(ip, Token)
		if err != nil {
			return "", err
		}
		if bool {
			return claims["uuid"].(string), nil
		}
		return "", errors.New("ip check Not pass")
	}
	return "", errors.New("claims exp, iat, nbf Not pass or ok")
}

// parseInterface ...
func parseInterface(token *jwt.Token) (interface{}, error) {
	return []byte(Aeskey), nil
}
