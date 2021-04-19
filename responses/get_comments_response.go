package responses

import (
	"time"
	"xendit-takehome/github/entities"
)

type GetCommentResponse struct {
	Id        int       `json:"id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy int       `json:"createdBy"`
}

type GetCommentsResponse []GetCommentResponse

func CreateGetCommentsResponse(entities []entities.Comment) GetCommentsResponse {
	response := GetCommentsResponse{}
	for _, entity := range entities {
		response = append(response, GetCommentResponse{
			Id:        entity.ID,
			Comment:   entity.Comment,
			CreatedAt: entity.CreatedAt,
			CreatedBy: entity.CreatedBy,
		})
	}
	return response
}
