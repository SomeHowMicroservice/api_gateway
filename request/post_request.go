package request

type CreateTopicRequest struct {
	Name string  `json:"name" binding:"required,max=50"`
	Slug *string `json:"slug" binding:"omitempty,max=50"`
}

type UpdateTopic struct {
	Name string `json:"name" binding:"required,max=50"`
	Slug string `json:"slug" binding:"required,max=50"`
}

type CreatePostRequest struct {
	Title       string `json:"title" binding:"required,min=1"`
	Content     string `json:"content" binding:"required,min=1"`
	IsPublished *bool  `json:"is_published" binding:"required"`
	TopicID     string `json:"topic_id" binding:"required,uuid4"`
}

type PostPaginationQuery struct {
	Page        uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit       uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort        string `form:"sort" json:"sort"`
	Order       string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	IsPublished *bool  `form:"is_active" json:"is_published"`
	Search      string `form:"search" json:"search"`
	TopicID     string `form:"topic_id" json:"topic_id" binding:"omitempty,uuid4"`
}
