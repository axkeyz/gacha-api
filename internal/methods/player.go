package methods

//============================= Currency ===========================//
// Player is the player (gamer) API.
//
// This is directly related to a player's identity in-game.
//
// This is directly mapped to the player table.
//===============================================================//
type Player struct {
	ID         int    `json:",omitempty"`
	Username   string `json:",omitempty"`
	Email      string `json:",omitempty"`
	Verified   bool   `json:",omitempty"`
	VerifiedAt string `json:",omitempty"`
	CreatedAt  string `json:",omitempty"`
	UpdatedAt  string `json:",omitempty"`
	DisabledAt string `json:",omitempty"`
}
