package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
// Handlers
// -----------------------------------------------------------------------------

func createMatch(c *gin.Context) {
    matchID := uuid.NewString()
    playerID := uuid.NewString()

    tok, err := newToken(matchID, playerID, "host")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    inviteURL := c.Request.Host + "/join?invite=" + tok

    // persistence layer for match header would go here

    c.JSON(http.StatusCreated, gin.H{
        "matchId": matchID,
        "playerId": playerID,
        "token":   tok,
        "invite":  inviteURL,
    })
}

func joinMatch(c *gin.Context) {
    inviteToken := c.Query("invite")
    if inviteToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing invite token"})
        return
    }

    claims, err := parseToken(inviteToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid invite token"})
        return
    }

    matchID := claims.Mid
    newPlayerID := uuid.NewString()

    playerTok, err := newToken(matchID, newPlayerID, "guest")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "matchId":  matchID,
        "playerId": newPlayerID,
        "token":    playerTok,
    })
}

// -----------------------------------------------------------------------------
// main router / CORS glue
// -----------------------------------------------------------------------------

func main() {
    r := gin.Default()

    // ---- CORS --------------------------------------------------------------
    config := cors.Config{
        AllowOrigins: []string{
            "http://localhost:5173", "http://localhost:3000",
            "http://127.0.0.1:5173", "http://127.0.0.1:3000",
        },
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }
    r.Use(cors.New(config))

    // ---- routes ------------------------------------------------------------
    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, Gin!")
    })

    r.POST("/matches", createMatch) // host creates a match
    r.GET("/join", joinMatch)       // guest accepts invite

    log.Println("listening on :8080 …")
    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
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
