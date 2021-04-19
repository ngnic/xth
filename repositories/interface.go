package repositories

import "xendit-takehome/github/entities"

type OrganisationRepository interface {
	GetComments(orgName string) ([]entities.Comment, error)
	AddComment(orgName string, username string, comment string) error
	DeleteComments(orgName string) error
	GetMembers(orgName string) ([]entities.Member, error)
}

type UserRepository interface {
	GetUser(apiKey string) (entities.UserOrganisation, error)
}
