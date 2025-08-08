// backend/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"

	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"server/internal/auth"
	"server/internal/database/domain"
	dbmigrate "server/internal/database/migrate"
	postgresrepo "server/internal/database/postgres"
)

// ----------------------------------------------------------------------------
// App – holds all shared dependencies for handlers
// ----------------------------------------------------------------------------
type App struct {
	MatchRepo domain.MatchRepository
}

// ----------------------------------------------------------------------------
// Handlers (now methods on *App)
// ----------------------------------------------------------------------------

func (a *App) createMatch(c *gin.Context) {
	matchID  := uuid.NewString()
	playerID := uuid.NewString()

	// Build initial domain object
	now := time.Now()
	match := &domain.WordleMatch{
		ID:        matchID,
		Status:    domain.MatchWaiting,
		CreatedAt: now,
		Players: [2]domain.Player{
			{ID: playerID, Connected: true},
		},
	}

	// Persist
	if err := a.MatchRepo.Create(c.Request.Context(), match); err != nil {
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
    // 1. Connect to PostgreSQL via GORM
    dsn := os.Getenv("DATABASE_URL") // e.g. postgres://user:pass@localhost:5432/wordle
    gormDB, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("opening DB: %v", err)
    }
    // Configure pool and ping
    sqlDB, err := gormDB.DB()
    if err != nil {
        log.Fatalf("acquire sql DB: %v", err)
    }
    sqlDB.SetMaxOpenConns(10)
    sqlDB.SetConnMaxIdleTime(5 * time.Minute)
    if err := sqlDB.PingContext(context.Background()); err != nil {
        log.Fatalf("ping DB: %v", err)
    }

    // 2. Run migrations
    if err := dbmigrate.Run(context.Background(), gormDB); err != nil {
        log.Fatalf("migrate: %v", err)
    }

    // 3. Create the repository
    repo := postgresrepo.NewMatchRepository(gormDB)

    // 4. Build the App with shared deps
	app := &App{MatchRepo: repo}

    // 5. Router setup
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

    // 6. Routes – bind methods
	r.POST("/matches", app.createMatch)
	r.GET("/join",   app.joinMatch)

	log.Println("listening on :8080 …")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}