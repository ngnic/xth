package repositories_test

import (
	"testing"
	"xendit-takehome/github/repositories"
	"xendit-takehome/github/testing_utils"
	"xendit-takehome/github/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestGetUser(t *testing.T) {
	db := testing_utils.GetDBHandle()
	defer testing_utils.CleanupTables(db)

	userRepo := repositories.NewUserDBRepository(db)

	genericPasswordHash, _ := utils.HashPassword("A")
	apiKey := uuid.NewString()

	//given
	orgIds := []int{}
	rows, _ := db.Query("insert into organisations(name) values ('a'), ('b') returning id")
	for rows.Next() {
		var id int
		rows.Scan(&id)
		orgIds = append(orgIds, id)
	}

	userIds := []int{}
	userRows, _ := db.NamedQuery(`
		insert into 
			users(username, password_hash, api_key) 
		values 
			(:userA, :passwordHash, :targetKey),
			(:userB, :passwordHash, :randomKey)
		returning id
	`, map[string]interface{}{
		"passwordHash": genericPasswordHash,
		"userA":        "A",
		"userB":        "B",
		"randomKey":    uuid.NewString(),
		"targetKey":    apiKey,
	})
	for userRows.Next() {
		var id int
		userRows.Scan(&id)
		userIds = append(userIds, id)
	}

	db.NamedExec(`
		insert into
			user_organisations(user_id, organisation_id)
		values
			(:userA, :orgA),
			(:userA, :orgB),
			(:userB, :orgB)
	`, map[string]interface{}{
		"userA": userIds[0],
		"userB": userIds[1],
		"orgA":  orgIds[0],
		"orgB":  orgIds[1],
	})

	//when
	user, _ := userRepo.GetUser(apiKey)

	//then
	assert.Equal(t, "A", user.Username)
	assert.Equal(t, userIds[0], user.ID)
	assert.ElementsMatch(t, []string{"a", "b"}, user.OrganisationNames)
}
