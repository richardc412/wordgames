// backend/internal/database/postgres/match_repository.go
// GORM-based implementation of the MatchRepository storing the domain object
// as a single JSONB column with row-level locking for atomic updates.
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"server/internal/database/domain"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNotFound = errors.New("match not found")

// matchRow is the GORM model for the matches table.
type matchRow struct {
    ID        string         `gorm:"primaryKey;type:text"`
    Data      datatypes.JSON `gorm:"type:jsonb;not null"`
    UpdatedAt time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

func (matchRow) TableName() string { return "matches" }

// MatchRepository is a PostgreSQL-backed implementation of domain.MatchRepository using GORM.
type MatchRepository struct {
    db *gorm.DB
}

// NewMatchRepository returns a repo backed by the given *gorm.DB.
func NewMatchRepository(db *gorm.DB) *MatchRepository {
    return &MatchRepository{db: db}
}

// fetchLocked loads and locks a row within the provided transaction.
func fetchLocked(ctx context.Context, tx *gorm.DB, id string) (*domain.WordleMatch, error) {
    var row matchRow
    if err := tx.WithContext(ctx).
        Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&row, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrNotFound
        }
        return nil, err
    }
    var m domain.WordleMatch
    if err := json.Unmarshal(row.Data, &m); err != nil {
        return nil, err
    }
    return &m, nil
}

// Create inserts a fresh match in the waiting state.
func (r *MatchRepository) Create(ctx context.Context, m *domain.WordleMatch) error {
    payload, err := json.Marshal(m)
    if err != nil {
        return err
    }
    row := matchRow{
        ID:        m.ID,
        Data:      datatypes.JSON(payload),
        UpdatedAt: time.Now(),
    }
    if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
        return err
    }
    return nil
}

// Get retrieves the match and returns a deep copy.
func (r *MatchRepository) Get(ctx context.Context, id string) (*domain.WordleMatch, error) {
    var row matchRow
    if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrNotFound
        }
        return nil, err
    }
    var m domain.WordleMatch
    if err := json.Unmarshal(row.Data, &m); err != nil {
        return nil, err
    }
    return &m, nil
}

// Start sets a match to in_progress and records the start timestamp.
func (r *MatchRepository) Start(ctx context.Context, id string, now time.Time) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        m, err := fetchLocked(ctx, tx, id)
        if err != nil {
            return err
        }
        m.Status = domain.MatchInProgress
        m.StartedAt = &now

        payload, err := json.Marshal(m)
        if err != nil {
            return err
        }
        if err := tx.Model(&matchRow{}).
            Where("id = ?", id).
            Updates(map[string]any{
                "data":       datatypes.JSON(payload),
                "updated_at": now,
            }).Error; err != nil {
            return err
        }
        return nil
    })
}

// AddGuess appends a guess for the specified player and updates their remaining clock.
func (r *MatchRepository) AddGuess(ctx context.Context, matchID string, playerID string, guess domain.Guess, remainingMs int64) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        m, err := fetchLocked(ctx, tx, matchID)
        if err != nil {
            return err
        }
        playerIndex := -1
        for i := range m.Players {
            if m.Players[i].ID == playerID {
                playerIndex = i
                break
            }
        }
        if playerIndex == -1 {
            return errors.New("player not part of match")
        }
        m.Players[playerIndex].Guesses = append(m.Players[playerIndex].Guesses, guess)
        m.Players[playerIndex].RemainingMs = remainingMs

        payload, err := json.Marshal(m)
        if err != nil {
            return err
        }
        if err := tx.Model(&matchRow{}).
            Where("id = ?", matchID).
            Updates(map[string]any{
                "data":       datatypes.JSON(payload),
                "updated_at": time.Now(),
            }).Error; err != nil {
            return err
        }
        return nil
    })
}

// SetPlayerConnected toggles the connected flag for a player.
func (r *MatchRepository) SetPlayerConnected(ctx context.Context, matchID, playerID string, connected bool) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        m, err := fetchLocked(ctx, tx, matchID)
        if err != nil {
            return err
        }
        playerIndex := -1
        for i := range m.Players {
            if m.Players[i].ID == playerID {
                playerIndex = i
                break
            }
        }
        if playerIndex == -1 {
            return errors.New("player not part of match")
        }
        m.Players[playerIndex].Connected = connected

        payload, err := json.Marshal(m)
        if err != nil {
            return err
        }
        if err := tx.Model(&matchRow{}).
            Where("id = ?", matchID).
            Updates(map[string]any{
                "data":       datatypes.JSON(payload),
                "updated_at": time.Now(),
            }).Error; err != nil {
            return err
        }
        return nil
    })
}

// Finish marks the match as finished with the given end state and winner.
func (r *MatchRepository) Finish(ctx context.Context, matchID, playerID string, endState domain.FinishedState, now time.Time) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        m, err := fetchLocked(ctx, tx, matchID)
        if err != nil {
            return err
        }
        m.Status = domain.MatchFinished
        m.EndedAt = &now
        if playerID != "" {
            for i := range m.Players {
                if m.Players[i].ID == playerID {
                    m.Players[i].FinishedState = endState
                }
            }
        }
        payload, err := json.Marshal(m)
        if err != nil {
            return err
        }
        if err := tx.Model(&matchRow{}).
            Where("id = ?", matchID).
            Updates(map[string]any{
                "data":       datatypes.JSON(payload),
                "updated_at": now,
            }).Error; err != nil {
            return err
        }
        return nil
    })
}

// Delete removes a match entirely.
func (r *MatchRepository) Delete(ctx context.Context, id string) error {
    if err := r.db.WithContext(ctx).Delete(&matchRow{}, "id = ?", id).Error; err != nil {
        return err
    }
    return nil
}


