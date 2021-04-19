package requests

type PostCommentRequest struct {
	Comment string `form:"comment" json:"comment" binding:"required"`
}
