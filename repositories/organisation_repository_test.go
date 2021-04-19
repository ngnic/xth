package repositories_test

import (
	"testing"
	"xendit-takehome/github/repositories"
	"xendit-takehome/github/testing_utils"
	"xendit-takehome/github/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetComments(t *testing.T) {
	db := testing_utils.GetDBHandle()
	defer testing_utils.CleanupTables(db)

	orgRepo := repositories.NewOrganisationDBRepository(db)

	genericPasswordHash, _ := utils.HashPassword("A")

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
			(:userA, :genericHash, :keyA), 
			(:userB, :genericHash, :keyB)
		returning id
	`, map[string]interface{}{
		"userA":       "A",
		"userB":       "B",
		"genericHash": genericPasswordHash,
		"keyA":        uuid.NewString(),
		"keyB":        uuid.NewString(),
	})
	for userRows.Next() {
		var id int
		userRows.Scan(&id)
		userIds = append(userIds, id)
	}

	// insert 4 rows
	// row 3 has a different organisation
	// row 4 is deleted
	db.NamedExec(`
		insert into 
			organisation_comments(comment, created_by, organisation_id, deleted_at)
		values
			('1', :userA, :orgA, null),
			('2', :userB, :orgA, null),
			('3', :userA, :orgB, null),
			('4', :userA, :orgA, now())
	`, map[string]interface{}{
		"userA": userIds[0],
		"userB": userIds[1],
		"orgA":  orgIds[0],
		"orgB":  orgIds[1],
	})

	//when
	comments, _ := orgRepo.GetComments("a")

	//then
	assert.Len(t, comments, 2)
	assert.True(t, comments[1].ID > comments[0].ID)
	assert.True(t, comments[0].OrganisationId == comments[1].OrganisationId)

	assert.Equal(t, "1", comments[0].Comment)
	assert.Equal(t, userIds[0], comments[0].CreatedBy)

	assert.Equal(t, "2", comments[1].Comment)
	assert.Equal(t, userIds[1], comments[1].CreatedBy)

	db.MustExec("truncate organisation_comments cascade")
	db.MustExec("truncate users cascade")
	db.MustExec("truncate organisations cascade")
}

func TestDeleteComments(t *testing.T) {
	db := testing_utils.GetDBHandle()
	defer testing_utils.CleanupTables(db)

	orgRepo := repositories.NewOrganisationDBRepository(db)

	genericPasswordHash, _ := utils.HashPassword("A")

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
			(:userA, :genericHash, :keyA), 
			(:userB, :genericHash, :keyB)
		returning id
	`, map[string]interface{}{
		"userA":       "A",
		"userB":       "B",
		"genericHash": genericPasswordHash,
		"keyA":        uuid.NewString(),
		"keyB":        uuid.NewString(),
	})
	for userRows.Next() {
		var id int
		userRows.Scan(&id)
		userIds = append(userIds, id)
	}

	// insert 4 rows
	// row 3 has a different organisation
	// row 4 is deleted
	db.NamedExec(`
		insert into 
			organisation_comments(comment, created_by, organisation_id, deleted_at)
		values
			('1', :userA, :orgA, null),
			('2', :userB, :orgA, null),
			('3', :userA, :orgB, null),
			('4', :userA, :orgA, now())
	`, map[string]interface{}{
		"userA": userIds[0],
		"userB": userIds[1],
		"orgA":  orgIds[0],
		"orgB":  orgIds[1],
	})

	/*
		// insert 4 rows
		// row 3 has a different organisation
		db.MustExec(`
			insert into
				organisation_comments(comment, created_by, organisation_id, deleted_at)
			values
				($1, $2, $3, $4),
				($5, $6, $7, $8),
				($9, $10, $11, $12),
				($13, $14, $15, $16)
		`, "1", userIds[0], orgIds[0], nil, "2", userIds[1], orgIds[0], nil, "3", userIds[0], orgIds[1], nil, "4", userIds[0], orgIds[0], time.Now())
	*/
	//when
	orgRepo.DeleteComments("a")

	//then
	comments, _ := orgRepo.GetComments("a")
	assert.Len(t, comments, 0)

	//rows deleted before should be left untouched
	var uniqueTimestamps int
	db.QueryRowx(`
		select count(distinct deleted_at) as unique_timestamp_count from organisation_comments where organisation_id = $1
	`, orgIds[0]).Scan(&uniqueTimestamps)
	assert.Equal(t, 2, uniqueTimestamps)

	//organisation b's comment should not be deleted
	comments, _ = orgRepo.GetComments("b")
	assert.Len(t, comments, 1)
}

func TestAddComment(t *testing.T) {
	db := testing_utils.GetDBHandle()
	defer testing_utils.CleanupTables(db)

	orgRepo := repositories.NewOrganisationDBRepository(db)

	genericPasswordHash, _ := utils.HashPassword("A")

	//given
	var orgId int
	db.QueryRowx("insert into organisations(name) values ('test org')").Scan(&orgId)

	var userId int
	db.QueryRowx(`
		insert into 
			users(username, password_hash, api_key) values ($1, $2, $3) returning id
	`, "A", genericPasswordHash, "123456").Scan(&userId)

	//when
	orgRepo.AddComment("test org", "A", "test-comment")

	//then
	comments, _ := orgRepo.GetComments("test org")
	assert.Len(t, comments, 1)
	assert.Equal(t, "test-comment", comments[0].Comment)
	assert.Equal(t, userId, comments[0].CreatedBy)
	assert.Equal(t, orgId, comments[0].OrganisationId)
}

func TestGetMembers(t *testing.T) {
	db := testing_utils.GetDBHandle()
	defer testing_utils.CleanupTables(db)

	orgRepo := repositories.NewOrganisationDBRepository(db)

	genericPasswordHash, _ := utils.HashPassword("A")

	//given
	orgIds := []int{}
	orgRows, _ := db.Query("insert into organisations(name) values ('org A'), ('org B') returning id")
	for orgRows.Next() {
		var id int
		orgRows.Scan(&id)
		orgIds = append(orgIds, id)
	}

	userIds := []int{}
	userRows, _ := db.NamedQuery(`
		insert into 
			users(username, password_hash, api_key) 
		values 
			(:userA, :passwordHash, :keyA),
			(:userB, :passwordHash, :keyB),
			(:userC, :passwordHash, :keyC),
			(:userD, :passwordHash, :keyD)
		returning id
	`, map[string]interface{}{
		"passwordHash": genericPasswordHash,
		"userA":        "A",
		"userB":        "B",
		"userC":        "C",
		"userD":        "D",
		"keyA":         uuid.NewString(),
		"keyB":         uuid.NewString(),
		"keyC":         uuid.NewString(),
		"keyD":         uuid.NewString(),
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
			(:userB, :orgA),
			(:userC, :orgB),
			(:userD, :orgB)
	`, map[string]interface{}{
		"userA": userIds[0],
		"userB": userIds[1],
		"userC": userIds[2],
		"userD": userIds[3],
		"orgA":  orgIds[0],
		"orgB":  orgIds[1],
	})

	db.NamedExec(`
		insert into
			user_followers(followee_id, follower_id)
		values
			(:userA, :userB),
			(:userB, :userA),
			(:userB, :userD),
			(:userC, :userA),
			(:userD, :userB)
	`, map[string]interface{}{
		"userA": userIds[0],
		"userB": userIds[1],
		"userC": userIds[2],
		"userD": userIds[3],
	})
	//when
	members, _ := orgRepo.GetMembers("org A")

	//then
	assert.Len(t, members, 2)
	assert.Equal(t, userIds[1], members[0].ID)
	assert.Equal(t, userIds[0], members[1].ID)
	assert.True(t, members[0].FollowedBy > members[1].FollowedBy)
	assert.Equal(t, 2, members[0].FollowedBy)
	assert.Equal(t, 2, members[0].Following)

	assert.Equal(t, 1, members[1].FollowedBy)
	assert.Equal(t, 2, members[1].Following)
}
