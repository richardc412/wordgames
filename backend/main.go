// backend/main.go
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq" // <-- driver

	"server/internal/auth"
	"server/internal/database/domain"
	"server/internal/database/postgres"
)

// ----------------------------------------------------------------------------
// App – holds all shared dependencies for handlers
// ----------------------------------------------------------------------------
type App struct {
	Repo domain.MatchRepository
}

// ----------------------------------------------------------------------------
// Handlers (now methods on *App)
// ----------------------------------------------------------------------------

func (a *App) createMatch(c *gin.Context) {
	matchID  := uuid.NewString()
	playerID := uuid.NewString()

	// Build initial domain object
	now := time.Now()
	m := &domain.WordleMatch{
		ID:        matchID,
		Status:    domain.MatchWaiting,
		CreatedAt: now,
		Players: [2]domain.Player{
			{ID: playerID, Connected: true},
		},
	}

	// Persist
	if err := a.Repo.Create(c.Request.Context(), m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auth token the client will use for future requests / websocket
	tok, _ := auth.NewToken(matchID, playerID, "host")

	c.JSON(http.StatusCreated, gin.H{
		"matchId":  matchID,
		"playerId": playerID,
		"token":    tok,
	})
}

func (a *App) joinMatch(c *gin.Context) {
	// … your existing code …
}

// ----------------------------------------------------------------------------
// main – initialise DB, repo, router, etc.
// ----------------------------------------------------------------------------
func main() {
	// 1. Connect to PostgreSQL
	dsn := os.Getenv("DATABASE_URL") // e.g. postgres://user:pass@localhost:5432/wordle
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("opening DB: %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Optional: ping on startup to fail fast
	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("ping DB: %v", err)
	}

	// 2. Create the repository
	repo := postgres.NewMatchRepository(db)

	// 3. Build the App with shared deps
	app := &App{Repo: repo}

	// 4. Router setup
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 5. Routes – bind methods
	r.POST("/matches", app.createMatch)
	r.GET("/join",   app.joinMatch)

	log.Println("listening on :8080 …")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}