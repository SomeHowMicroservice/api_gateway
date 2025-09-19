package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/gateway/common"
	postpb "github.com/SomeHowMicroservice/gateway/protobuf/post"
	userpb "github.com/SomeHowMicroservice/gateway/protobuf/user"
	"github.com/SomeHowMicroservice/gateway/request"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postClient postpb.PostServiceClient
}

func NewPostHandler(postClient postpb.PostServiceClient) *PostHandler {
	return &PostHandler{postClient}
}

func (h *PostHandler) CreateTopic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	var req request.CreateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	var slug *string
	if req.Slug != nil {
		slug = req.Slug
	}

	res, err := h.postClient.CreateTopic(ctx, &postpb.CreateTopicRequest{
		Name:   req.Name,
		Slug:   slug,
		UserId: user.Id,
	})
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusCreated, "Tạo chủ đề bài viết thành công", res)
}

func (h *PostHandler) GetAllTopicsAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.postClient.GetAllTopicsAdmin(ctx, &postpb.GetAllRequest{})
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Lấy danh sách chủ đề bài viết thành công", res)
}

func (h *PostHandler) GetAllTopics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.postClient.GetAllTopics(ctx, &postpb.GetAllRequest{})
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Lấy danh sách chủ đề bài viết thành công", res)
}

func (h *PostHandler) GetDeletedTopics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.postClient.GetDeletedTopics(ctx, &postpb.GetAllRequest{})
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả chủ đề đã xóa thành công", res)
}

func (h *PostHandler) UpdateTopic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	var req request.UpdateTopic
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	topicID := c.Param("id")

	if _, err := h.postClient.UpdateTopic(ctx, &postpb.UpdateTopicRequest{
		Id:     topicID,
		Name:   req.Name,
		Slug:   req.Slug,
		UserId: user.Id,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Cập nhật chủ đề bài viết thành công", nil)
}

func (h *PostHandler) DeleteTopic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	topicID := c.Param("id")

	if _, err := h.postClient.DeleteTopic(ctx, &postpb.DeleteOneRequest{
		Id:     topicID,
		UserId: user.Id,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Chuyển chủ đề bài viết vào thùng rác thành công", nil)
}

func (h *PostHandler) DeleteTopics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	var req request.DeleteManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.postClient.DeleteTopics(ctx, &postpb.DeleteManyRequest{
		Ids:    req.IDs,
		UserId: user.Id,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Chuyển danh sách chủ đề vào thùng rác thành công", nil)
}

func (h *PostHandler) RestoreTopic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	topicID := c.Param("id")

	if _, err := h.postClient.RestoreTopic(ctx, &postpb.RestoreOneRequest{
		Id:     topicID,
		UserId: user.Id,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Khôi phục chủ đề bài viết thành công", nil)
}

func (h *PostHandler) RestoreTopics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	var req request.RestoreManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.postClient.RestoreTopics(ctx, &postpb.RestoreManyRequest{
		Ids:    req.IDs,
		UserId: user.Id,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Khôi phục danh sách chủ đề thành công", nil)
}

func (h *PostHandler) PermanentlyDeleteTopic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	topicID := c.Param("id")

	if _, err := h.postClient.PermanentlyDeleteTopic(ctx, &postpb.PermanentlyDeleteOneRequest{
		Id: topicID,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Xóa chủ đề bài viết thành công", nil)
}

func (h *PostHandler) PermanentlyDeleteTopics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.DeleteManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.postClient.PermanentlyDeleteTopics(ctx, &postpb.PermanentlyDeleteManyRequest{
		Ids: req.IDs,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Xóa danh sách chủ đề thành công", nil)
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	var req request.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.postClient.CreatePost(ctx, &postpb.CreatePostRequest{
		Title:       req.Title,
		Content:     req.Content,
		IsPublished: *req.IsPublished,
		TopicId:     req.TopicID,
		UserId:      user.Id,
	})
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusCreated, "Tạo bài viết thành công", res)
}

func (h *PostHandler) GetAllPostsAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var query request.PostPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.postClient.GetAllPostsAdmin(ctx, &postpb.GetAllPostsAdminRequest{
		Page:        query.Page,
		Limit:       query.Limit,
		Sort:        query.Sort,
		Order:       query.Order,
		Search:      query.Search,
		TopicId:     query.TopicID,
		IsPublished: query.IsPublished,
	})
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả bài viết thành công", res)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	postId := c.Param("id")

	res, err := h.postClient.GetPostById(ctx, &postpb.GetOneRequest{
		Id: postId,
	})
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Lấy chi tiết bài viết thành công", gin.H{
		"post": res,
	})
}

func (h *PostHandler) GetPostContentByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	postId := c.Param("id")

	res, err := h.postClient.GetPostContentById(ctx, &postpb.GetOneRequest{
		Id: postId,
	})
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Lấy nội dung bài viết thành công", res)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	var req request.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	postID := c.Param("id")

	var title, content, topicID string
	var isPublished bool
	if req.Title != nil {
		title = *req.Title
	}
	if req.Content != nil {
		content = *req.Content
	}
	if req.TopicID != nil {
		topicID = *req.TopicID
	}
	if req.IsPublished != nil {
		isPublished = *req.IsPublished
	}

	res, err := h.postClient.UpdatePost(ctx, &postpb.UpdatePostRequest{
		Id:          postID,
		Title:       &title,
		Content:     &content,
		TopicId:     &topicID,
		IsPublished: &isPublished,
		UserId:      user.Id,
	}) 
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Chỉnh sửa bài viết thành công", gin.H{
		"post": res,
	})
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	postID := c.Param("id")

	if _, err := h.postClient.DeletePost(ctx, &postpb.DeleteOneRequest{
		Id:     postID,
		UserId: user.Id,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Chuyển bài viết vào thùng rác thành công", nil)
}

func (h *PostHandler) DeletePosts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	var req request.DeleteManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.postClient.DeletePosts(ctx, &postpb.DeleteManyRequest{
		Ids:    req.IDs,
		UserId: user.Id,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Chuyển danh sách bài viết vào thùng rác thành công", nil)
}

func (h *PostHandler) RestorePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	postID := c.Param("id")

	if _, err := h.postClient.RestorePost(ctx, &postpb.RestoreOneRequest{
		Id:     postID,
		UserId: user.Id,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Khôi phục bài viết vào thành công", nil)
}

func (h *PostHandler) RestorePosts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	var req request.RestoreManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.postClient.RestorePosts(ctx, &postpb.RestoreManyRequest{
		Ids:    req.IDs,
		UserId: user.Id,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Khôi phục danh sách bài viết thành công", nil)
}

func (h *PostHandler) PermanentlyDeletePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	postID := c.Param("id")

	if _, err := h.postClient.PermanentlyDeletePost(ctx, &postpb.PermanentlyDeleteOneRequest{
		Id: postID,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Xóa bài viết thành công", nil)
}

func (h *PostHandler) PermanentlyDeletePosts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.DeleteManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.postClient.PermanentlyDeletePosts(ctx, &postpb.PermanentlyDeleteManyRequest{
		Ids: req.IDs,
	}); common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Xóa danh sách bài viết thành công", nil)
}

func (h *PostHandler) GetDeletedPosts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var query request.PostPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.postClient.GetDeletedPosts(ctx, &postpb.GetAllPostsAdminRequest{
		Page:        query.Page,
		Limit:       query.Limit,
		Sort:        query.Sort,
		Order:       query.Order,
		Search:      query.Search,
		TopicId:     query.TopicID,
		IsPublished: query.IsPublished,
	})
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả bài viết đã xóa thành công", res)
}

func (h *PostHandler) GetDeletedPostByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	postID := c.Param("id")

	res, err := h.postClient.GetDeletedPostById(ctx, &postpb.GetOneRequest{
		Id: postID,
	})
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Lấy chi tiết bài viết thành công", gin.H{
		"post": res,
	})
}
