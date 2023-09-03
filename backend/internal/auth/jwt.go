package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthToken struct {
	AuthTime float64                `json:"auth_time"`
	Issuer   string                 `json:"iss"`
	Audience string                 `json:"aud"`
	Expires  float64                `json:"exp"`
	IssuedAt float64                `json:"iat"`
	Subject  string                 `json:"sub,omitempty"`
	UID      string                 `json:"uid,omitempty"`
	Claims   map[string]interface{} `json:"-"`
}

var keyRetreiver = NewGoogleKeyRetreiver()

func WarmCache(ctx context.Context) {
	keyRetreiver.GetKeys(ctx)
}

func ValidateToken(ctx context.Context, audience string, tokenStr string) (*AuthToken, error) {
	keys, err := keyRetreiver.GetKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting keys: %w", err)
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"]
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}
		kidStr, ok := kid.(string)
		if !ok {
			return nil, fmt.Errorf("kid was not a string")
		}
		key, ok := keys[kidStr]
		if !ok {
			return nil, fmt.Errorf("key not found for kid %v", kidStr)
		}
		block, _ := pem.Decode([]byte(key))
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate")
		}
		publicKey := cert.PublicKey.(*rsa.PublicKey)
		return publicKey, nil
	}, jwt.WithValidMethods([]string{"RS256"}), jwt.WithAudience(audience), jwt.WithIssuer("https://securetoken.google.com/"+audience))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token not valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to get token map claims")
	}
	authTimeIface, ok := claims["auth_time"]
	if !ok {
		return nil, fmt.Errorf("claim auth_time not found")
	}
	authTime, ok := authTimeIface.(float64)
	if !ok {
		return nil, fmt.Errorf("claim auth_time has invalid type")
	}
	// auth_time must be in the past
	now := float64(time.Now().UTC().Unix())
	if authTime > now {
		return nil, fmt.Errorf("auth_time must be in the past. auth_time=%f now=%f", authTime, now)
	}
	authToken := &AuthToken{
		Claims: make(map[string]interface{}),
	}
	for k, v := range claims {
		switch k {
		case "auth_time":
			authToken.AuthTime = v.(float64)
		case "iss":
			authToken.Issuer = v.(string)
		case "aud":
			authToken.Audience = v.(string)
		case "exp":
			authToken.Expires = v.(float64)
		case "iat":
			authToken.IssuedAt = v.(float64)
		case "sub":
			authToken.Subject = v.(string)
		case "user_id":
			authToken.UID = v.(string)
		default:
			authToken.Claims[k] = v
		}
	}
	return authToken, nil
}

func TokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}
