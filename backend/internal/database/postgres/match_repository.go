// backend/internal/database/postgres/match_repository.go
// Package postgres provides a concrete implementation of the domain.MatchRepository
// interface using PostgreSQL and a single JSONB column to hold the entire
// WordleMatch document. Each write happens inside a transaction and locks the
// row (SELECT .. FOR UPDATE) to guarantee atomicity and consistency when the
// game is updated concurrently by two websocket connections.
//
// Table DDL that matches this repository:
//
//	CREATE TABLE matches (
//	    id         text PRIMARY KEY,
//	    data       jsonb        NOT NULL,
//	    updated_at timestamptz  NOT NULL DEFAULT now()
//	);
//
// A partial index on (data ->> 'status') could be useful for lobbies, e.g.:
//
//	CREATE INDEX matches_waiting_idx ON matches ((data ->> 'status'))
//	WHERE (data ->> 'status') = 'waiting';
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"server/internal/database/domain"
)

// Ensure we have the pq driver (or pgx) linked if the user imports it in main.
// _ "github.com/lib/pq"

var ErrNotFound = errors.New("match not found")

// MatchRepository is a PostgreSQL-backed implementation of domain.MatchRepository.
type MatchRepository struct {
    db *sql.DB
}

// NewMatchRepository returns a repo backed by the given *sql.DB.
func NewMatchRepository(db *sql.DB) *MatchRepository {
    return &MatchRepository{db: db}
}

// helper: fetch + lock a match row inside an open Tx.
func fetchMatchTx(ctx context.Context, tx *sql.Tx, id string) (*domain.WordleMatch, error) {
    var payload []byte
    err := tx.QueryRowContext(ctx, `SELECT data FROM matches WHERE id = $1 FOR UPDATE`, id).Scan(&payload)
    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    }
    if err != nil {
        return nil, err
    }

    var m domain.WordleMatch
    if err := json.Unmarshal(payload, &m); err != nil {
        return nil, err
    }
    return &m, nil
}

// Create inserts a fresh match in the waiting state.
func (r *MatchRepository) Create(ctx context.Context, m *domain.WordleMatch) error {
    bytes, err := json.Marshal(m)
    if err != nil {
        return err
    }
    _, err = r.db.ExecContext(ctx, `INSERT INTO matches (id, data) VALUES ($1, $2)`, m.ID, bytes)
    return err
}

// Get retrieves the match and returns a deep copy.
func (r *MatchRepository) Get(ctx context.Context, id string) (*domain.WordleMatch, error) {
    var payload []byte
    err := r.db.QueryRowContext(ctx, `SELECT data FROM matches WHERE id = $1`, id).Scan(&payload)
    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    }
    if err != nil {
        return nil, err
    }
    var m domain.WordleMatch
    if err := json.Unmarshal(payload, &m); err != nil {
        return nil, err
    }
    return &m, nil
}

// Start sets a match to in_progress and records the start timestamp.
func (r *MatchRepository) Start(ctx context.Context, id string, now time.Time) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() { _ = tx.Rollback() }()

    m, err := fetchMatchTx(ctx, tx, id)
    if err != nil {
        return err
    }

    // Update fields.
    m.Status = domain.MatchInProgress
    m.StartedAt = &now

    bytes, err := json.Marshal(m)
    if err != nil {
        return err
    }

    if _, err := tx.ExecContext(ctx, `UPDATE matches SET data = $1, updated_at = $2 WHERE id = $3`, bytes, now, id); err != nil {
        return err
    }
    return tx.Commit()
}

// AddGuess appends a guess for the specified player and updates their remaining clock.
func (r *MatchRepository) AddGuess(ctx context.Context, matchID string, playerID string, guess domain.Guess, remainingMs int64) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() { _ = tx.Rollback() }()

    m, err := fetchMatchTx(ctx, tx, matchID)
    if err != nil {
        return err
    }

    // Find player index.
    idx := -1
    for i := range m.Players {
        if m.Players[i].ID == playerID {
            idx = i
            break
        }
    }
    if idx == -1 {
        return errors.New("player not part of match")
    }

    // Mutate.
    m.Players[idx].Guesses = append(m.Players[idx].Guesses, guess)
    m.Players[idx].RemainingMs = remainingMs

    bytes, err := json.Marshal(m)
    if err != nil {
        return err
    }

    if _, err := tx.ExecContext(ctx, `UPDATE matches SET data = $1, updated_at = $2 WHERE id = $3`, bytes, time.Now(), matchID); err != nil {
        return err
    }
    return tx.Commit()
}

// SetPlayerConnected toggles the connected flag for a player.
func (r *MatchRepository) SetPlayerConnected(ctx context.Context, matchID, playerID string, connected bool) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() { _ = tx.Rollback() }()

    m, err := fetchMatchTx(ctx, tx, matchID)
    if err != nil {
        return err
    }

    idx := -1
    for i := range m.Players {
        if m.Players[i].ID == playerID {
            idx = i
            break
        }
    }
    if idx == -1 {
        return errors.New("player not part of match")
    }

    m.Players[idx].Connected = connected

    bytes, err := json.Marshal(m)
    if err != nil {
        return err
    }

    if _, err := tx.ExecContext(ctx, `UPDATE matches SET data = $1, updated_at = $2 WHERE id = $3`, bytes, time.Now(), matchID); err != nil {
        return err
    }
    return tx.Commit()
}

// Finish marks the match as finished with the given end state and winner.
func (r *MatchRepository) Finish(ctx context.Context, matchID, playerID string, endState domain.FinishedState, now time.Time) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() { _ = tx.Rollback() }()

    m, err := fetchMatchTx(ctx, tx, matchID)
    if err != nil {
        return err
    }

    m.Status = domain.MatchFinished
    m.EndedAt = &now

    if playerID != "" {
        for i := range m.Players {
            if m.Players[i].ID == playerID {
                m.Players[i].FinishedState = endState
            } else {
                // If other player is still none, you might decide what to do; leave unchanged.
            }
        }
    }

    bytes, err := json.Marshal(m)
    if err != nil {
        return err
    }

    if _, err := tx.ExecContext(ctx, `UPDATE matches SET data = $1, updated_at = $2 WHERE id = $3`, bytes, now, matchID); err != nil {
        return err
    }
    return tx.Commit()
}

// Delete removes a match entirely.
func (r *MatchRepository) Delete(ctx context.Context, id string) error {
    _, err := r.db.ExecContext(ctx, `DELETE FROM matches WHERE id = $1`, id)
    return err
}

