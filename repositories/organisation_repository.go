package repositories

import (
	"fmt"
	"xendit-takehome/github/entities"

	"github.com/jmoiron/sqlx"
)

type OrganisationDBRepository struct {
	db *sqlx.DB
}

func NewOrganisationDBRepository(db *sqlx.DB) OrganisationRepository {
	return &OrganisationDBRepository{
		db: db,
	}
}

func (repo *OrganisationDBRepository) GetComments(orgName string) ([]entities.Comment, error) {
	allComments := []entities.Comment{}
	err := repo.db.Select(&allComments, `
		select 
			org_comments.id,
			org_comments.comment,
			org_comments.created_at,
			org_comments.created_by
		from organisation_comments as org_comments
		inner join organisations as orgs on orgs.id = org_comments.organisation_id
		where org_comments.deleted_at is null and orgs.name = $1 order by org_comments.id`, orgName)
	return allComments, err
}

func (repo *OrganisationDBRepository) AddComment(orgName string, username string, comment string) error {
	userId := -1
	orgId := repo.getOrgId(orgName)
	repo.db.QueryRowx(`
		select id from users where username = $1 limit 1
	`, username).Scan(&userId)
	repo.db.Exec(
		`insert into organisation_comments(comment, organisation_id, created_by) values ($1, $2, $3)`,
		comment, orgId, userId)
	return nil
}

func (repo *OrganisationDBRepository) DeleteComments(orgName string) error {
	orgId := repo.getOrgId(orgName)
	repo.db.Exec(
		`update organisation_comments set deleted_at = now() where organisation_id = $1 and deleted_at is null`,
		orgId)
	return nil
}

func (repo *OrganisationDBRepository) GetMembers(orgName string) ([]entities.Member, error) {
	allMembers := []entities.Member{}
	err := repo.db.Select(&allMembers, `
	select 
		u.id,
		u.username,
		u.avatar,
		u.created_at,
		(select count(1) from user_followers as uf where uf.follower_id = u.id) as following,
		(select count(1) from user_followers as uf where uf.followee_id = u.id) as followed_by
	from users u 
	inner join user_organisations uo on uo.user_id = u.id
	inner join organisations og on og.id = uo.organisation_id
	where u.deleted_at is null and og.name = $1 order by followed_by desc
	`, orgName)
	return allMembers, err
}

func (repo *OrganisationDBRepository) getOrgId(orgName string) int {
	orgId := -1
	err := repo.db.QueryRow(`
		select id from organisations as orgs where orgs.name = $1 limit 1
	`, orgName).Scan(&orgId)
	fmt.Println(err)
	return orgId
}
