package entities

import "time"

type Comment struct {
	ID             int
	Comment        string
	CreatedAt      time.Time `db:"created_at"`
	CreatedBy      int       `db:"created_by"`
	OrganisationId int       `db:"organisation_id"`
}
