// backend/internal/database/domain/match_repository.go
package domain

import (
	"context"
	"time"
)

type MatchRepository interface {
	// Create an empty match in “waiting” state.
	Create(ctx context.Context, m *WordleMatch) error

	// Fetch the match by id – returns a deep copy so callers can mutate safely.
	Get(ctx context.Context, id string) (*WordleMatch, error)

	// Begin -> In-progress, sets Start time & initial clocks atomically.
	Start(ctx context.Context, id string, now time.Time) error

	// Append a guess and update player clock in one transaction.
	AddGuess(
		ctx context.Context,
		matchID string,
		playerID string,
		guess Guess,
		remainingMs int64,
	) error

	// Mark a player’s connection status (true on connect / false on disconnect).
	SetPlayerConnected(ctx context.Context, matchID, playerID string, connected bool) error

	// Mark match finished & winner in one shot.
	Finish(
		ctx context.Context,
		matchID string,
		playerID string,            // winner or empty for draw/timeout
		endState FinishedState,     // success, timeout, resign
		now time.Time,
	) error

	// Hard-delete if you need administrative cleanup.
	Delete(ctx context.Context, id string) error
}