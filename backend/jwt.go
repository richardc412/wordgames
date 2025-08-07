package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// -----------------------------------------------------------------------------
// JWT setup & helpers
// -----------------------------------------------------------------------------

// MatchClaims represents the custom + registered JWT claims we embed in every
// token.
//
//  mid → Match ID (UUID string)
//  pid → Player ID (UUID string)
//  rol → "host" | "guest"
//
// We embed jwt.RegisteredClaims so we also get exp, iat, nbf, jti, etc.
// -----------------------------------------------------------------------------

type MatchClaims struct {
	Mid string `json:"mid"`
	Pid string `json:"pid"`
	Rol string `json:"rol"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// .env file doesn't exist, which is fine for development
		log.Println("No .env file found, using system environment variables")
	}

	// Get JWT secret from environment
	secret := getenv("JWT_SECRET", "dev_secret_change_me")
	if secret == "dev_secret_change_me" {
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET in .env file for production.")
	}
	jwtSecret = []byte(secret)
}

// tokenTTL is the lifetime of a freshly‑minted game token.
const tokenTTL = 12 * time.Hour

func newToken(matchID, playerID, role string) (string, error) {
	now := time.Now()
	claims := MatchClaims{
		Mid: matchID,
		Pid: playerID,
		Rol: role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
			ID:        uuid.NewString(), // jti – optional replay defence
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func parseToken(tok string) (*MatchClaims, error) {
	parsed, err := jwt.ParseWithClaims(tok, &MatchClaims{}, func(_ *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*MatchClaims)
	if !ok || !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// -----------------------------------------------------------------------------
// helpers
// -----------------------------------------------------------------------------

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
