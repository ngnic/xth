package entities

import "github.com/lib/pq"

type UserOrganisation struct {
	ID                int
	Username          string         `db:"username"`
	OrganisationNames pq.StringArray `db:"organisation_names"`
}
