package utils

const (
	REPLAY_DIR       = "replays"
	REPLAY_SEPARATOR = '|'
)

var (
	DEBUG_VARS = []string{"round"} //example for debugging variables
)

type ReplayMove struct {
	OwnMove Direction `json:"N"`
	OppMove Direction `json:"P"`
	//additional fields may come here for debug purposes
	Others []string `json:"O"` //debugging variables
}

type ReplayFormat struct {
	FieldWidth  int `json:"W"`
	FieldHeight int `json:"H"`
	OwnX        int `json:"N"`
	OwnY        int `json:"M"`
	OppX        int `json:"P"`
	OppY        int `json:"Q"`
}
