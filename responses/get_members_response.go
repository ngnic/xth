package responses

import (
	"time"
	"xendit-takehome/github/entities"
)

type GetMemberResponse struct {
	Id         int       `json:"id"`
	Username   string    `json:"username"`
	Avatar     string    `json:"avatar"`
	CreatedAt  time.Time `json:"createdAt"`
	Following  int       `json:"following"`
	FollowedBy int       `json:"followedBy"`
}

type GetMembersResponse []GetMemberResponse

func CreateGetMembersResponse(entities []entities.Member) GetMembersResponse {
	response := GetMembersResponse{}
	for _, entity := range entities {
		response = append(response, GetMemberResponse{
			Id:         entity.ID,
			Username:   entity.Username,
			Avatar:     entity.Avatar,
			Following:  entity.Following,
			FollowedBy: entity.FollowedBy,
			CreatedAt:  entity.CreatedAt,
		})
	}
	return response
}
