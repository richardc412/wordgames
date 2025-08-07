package domain

type Match struct {
	ID                 string `json:"id"`
	PlayerOneID        string `json:"player_one_id"`
	PlayerTwoID        string `json:"player_two_id"`
	PlayerOneClock     int    `json:"player_one_clock"`
	PlayerTwoClock     int    `json:"player_two_clock"`
	PlayerOneConnected bool   `json:"player_one_connected"`
	PlayerTwoConnected bool   `json:"player_two_connected"`
	WordleWord         string `json:"wordle_word"`
}

type Action struct {
	ActionType ActionType `json:"action_type"`
	PlayerID   string     `json:"player_id"`
	Data       any        `json:"data"`
}

type ActionType int

const (
	ActionTypeCreateMatch ActionType = iota
	ActionTypeJoinMatch
	ActionTypeStartMatch
	ActionTypeMakeMove
	ActionTypeEndMatch
	ActionTypeDisconnect
	ActionTypeReconnect
	ActionTypeTimeout
	ActionTypeError
)

var actionTypeNames = map[ActionType]string{
	ActionTypeCreateMatch: "create_match",
	ActionTypeJoinMatch:   "join_match",
	ActionTypeStartMatch:  "start_match",
	ActionTypeMakeMove:    "make_move",
	ActionTypeEndMatch:    "end_match",
	ActionTypeDisconnect:  "disconnect",
	ActionTypeReconnect:   "reconnect",
	ActionTypeTimeout:     "timeout",
	ActionTypeError:       "error",
}