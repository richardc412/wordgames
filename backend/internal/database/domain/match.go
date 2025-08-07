// backend/internal/database/domain/match.go
package domain

import "time"

/*
A WordleMatch is the **only** record you need to fully restore a game.
Everything else (guesses, timers, connectivity) is embedded inside it.
Store it as a single JSON document or a relational row with
JSONB columns.
*/
type WordleMatch struct {
	ID        string            // primary key
	Status    MatchStatus       // waiting, in_progress, finished
	Word      string            // the solution
	Players   [2]Player        	// exactly two players
	CreatedAt time.Time   
	StartedAt *time.Time 
	EndedAt   *time.Time  
}

/*
State that changes per-player during the match.
Everything is local to this player so we can update it atomically.
*/
type Player struct {
	ID            string      
	RemainingMs   int64         // chess-clock style
	Connected     bool          
	Guesses       []Guess       // up to 6
	FinishedState FinishedState // none, success, timeout, resign
}

/*
Each guess carries the word and server-side evaluation so
clients don’t have to recompute it after a refresh.
*/
type Guess struct {
	Word      string      
	Evaluation [5]Letter   
	Timestamp time.Time   
}

/* --- Small enums ------------------------------------------------ */

type MatchStatus string
const (
	MatchWaiting    MatchStatus = "waiting"
	MatchInProgress MatchStatus = "in_progress"
	MatchFinished   MatchStatus = "finished"
)

type FinishedState string
const (
	FinishNone     FinishedState = "none"
	FinishSuccess  FinishedState = "success"   // guessed correctly
	FinishTimeout  FinishedState = "timeout"   // ran out of clock
	FinishResign   FinishedState = "resign"    // quit / disconnected too long
)

type Letter string // “G”, “Y”, “B”