package methods

import (
	"encoding/json"

	"github.com/axkeyz/gacha-api/internal/utils"
	"github.com/axkeyz/gacha-api/config"
)

type StaffLog struct {
	ID int `json:",omitempty"`
	StaffID int `json:",omitempty"`
	Staff Staff `json:",omitempty"`
	StaffActionID int `json:",omitempty"`
	StaffAction StaffAction `json:",omitempty"`
	Success bool `json:",omitempty"`
	Notes string `json:",omitempty"`
	IPAddress string `json:",omitempty"`
	CreatedAt string `json:",omitempty"`
}

func (log *StaffLog) Create(success bool, notes interface{}) {
	log.Success = success
	notesJSON, _ := json.Marshal(notes)
	log.Notes = string(notesJSON)

	if utils.HasNoEmptyParams(
		[]string{log.StaffAction.Name, log.Notes, log.IPAddress},
	) {
		// Setup database
		db := config.SetupDB()
		defer db.Close()

		// Setup query
		q := `INSERT into staff_log (staff_id, staff_action_id, success, notes, ip_address)
		SELECT $1, staff_action.id, $2, $3, $4 from staff_action WHERE 
		staff_action.name = $5;`

		db.Exec(
			q, log.StaffID, log.Success, log.Notes, log.IPAddress, log.StaffAction.Name,
		); 
	}
}