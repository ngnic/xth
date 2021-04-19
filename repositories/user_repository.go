package repositories

import (
	"xendit-takehome/github/entities"

	"github.com/jmoiron/sqlx"
)

type UserDBRepository struct {
	db *sqlx.DB
}

func NewUserDBRepository(db *sqlx.DB) UserRepository {
	return &UserDBRepository{
		db: db,
	}
}

func (repo *UserDBRepository) GetUser(apiKey string) (entities.UserOrganisation, error) {
	var userOrganisation entities.UserOrganisation
	err := repo.db.QueryRowx(`
		select 
			u.id,
			u.username,
			array_agg(o.name)::text[] as organisation_names
		from users u
		inner join user_organisations uo on uo.user_id = u.id
		inner join organisations o on o.id = uo.organisation_id
		where api_key = $1 group by u.id, u.username 
	`, apiKey).StructScan(&userOrganisation)
	return userOrganisation, err
}
