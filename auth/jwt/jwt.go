package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}
type Claims struct {
	UserId int64 `json:"userId"`
	Exp    int64 `json:"exp"`
}

const (
	secretKey = "1qaz2wsx3edc"
)

var (
	ErrTokenExpired = errors.New("token has expired")
	ErrInvalidToken = errors.New("invalid token")
)

func Generate(p Claims) string {
	const op = "auth.jwt.Generate"

	h := header{
		Algorithm: "SHA256",
		Type:      "JWT",
	}
	headerJSON, _ := json.Marshal(h)
	payloadJSON, _ := json.Marshal(p)

	encodedHeader := encodeSegment(headerJSON)
	encodedPayload := encodeSegment(payloadJSON)
	signature := createSignature(encodedHeader + "." + encodedPayload)

	return fmt.Sprintf("%s.%s.%s", encodedHeader, encodedPayload, signature)
}

func encodeSegment(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func createSignature(data string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return encodeSegment(h.Sum(nil))
}

func Parse(token string) (*Claims, error) {
	const op = "auth.jwt.Parse"

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	signature := createSignature(parts[0] + "." + parts[1])
	if signature != parts[2] {
		return nil, ErrInvalidToken
	}

	payload, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var claims Claims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if claims.Exp < time.Now().Unix() {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}
