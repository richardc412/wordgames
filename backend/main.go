package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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


