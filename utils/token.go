package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(ttl time.Duration, payload interface{}, privateKey []byte) (string, error) {
	// decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	// if err != nil {
	// 	return "", fmt.Errorf("could not decode key: %w", err)
	// }

	// myString := string(signBytes[:])
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)

	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}
func ValidateToken(token string, publicKey []byte) (interface{}, error) {
	// decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not decode: %w", err)
	// }

	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)

	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims["sub"], nil
}
