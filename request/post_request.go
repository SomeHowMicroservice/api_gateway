package request

import "mime/multipart"

type CreateTopicRequest struct {
	Name string  `json:"name" binding:"required,max=50"`
	Slug *string `json:"slug" binding:"omitempty,max=50"`
}

type UpdateTopic struct {
	Name string `json:"name" binding:"required,max=50"`
	Slug string `json:"slug" binding:"required,max=50"`
}

type CreatePostRequest struct {
	Title       string `json:"title" binding:"required,min=2"`
	Content     string `json:"content" binding:"required,min=1"`
	IsPublished *bool  `json:"is_published" binding:"required"`
	TopicID     string `json:"topic_id" binding:"required,uuid4"`
}

type CreatePostForm struct {
	Title       string                `form:"title" validate:"required,min=2"`
	Content     string                `form:"content" validate:"required,min=1"`
	IsPublished *bool                 `form:"is_published" validate:"required"`
	TopicID     string                `form:"topic_id" validate:"required,uuid4"`
	Images      []CreatePostImageForm `form:"images" validate:"required,dive"`
}

type CreatePostImageForm struct {
	IsThumbnail *bool                 `form:"is_thumbnail" validate:"required"`
	SortOrder   int                   `form:"sort_order" validate:"required,gt=0"`
	File        *multipart.FileHeader `form:"file" validate:"required"`
}
