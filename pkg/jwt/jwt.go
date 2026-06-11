package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Claims struct {
	Subject    string `json:"sub"`
	Expiration int64  `json:"exp"`
	IssuedAt   int64  `json:"iat,omitempty"`
	Type       string `json:"typ,omitempty"`
	Role       string `json:"role,omitempty"`
	UserID     int64  `json:"uid,omitempty"`
	TokenType  string `json:"token_type,omitempty"`
}

func ParseSubjectHS256(token, secret string) (string, error) {
	claims, err := ParseClaimsHS256(token, secret)
	if err != nil {
		return "", err
	}
	if claims.Subject == "" {
		return "", errors.New("token subject empty")
	}
	return claims.Subject, nil
}

func ParseClaimsHS256(token, secret string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, errors.New("invalid token format")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, errors.New("invalid token header")
	}

	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return Claims{}, errors.New("invalid token header json")
	}
	if header.Alg != "HS256" {
		return Claims{}, errors.New("unsupported jwt algorithm")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return Claims{}, errors.New("invalid token signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, errors.New("invalid token payload")
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return Claims{}, errors.New("invalid token payload json")
	}

	if claims.Expiration > 0 && time.Now().Unix() > claims.Expiration {
		return Claims{}, errors.New("token expired")
	}

	return claims, nil
}

func SignHS256(claims Claims, secret string) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)
	unsigned := headerEncoded + "." + payloadEncoded

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsigned))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return unsigned + "." + signature, nil
}
