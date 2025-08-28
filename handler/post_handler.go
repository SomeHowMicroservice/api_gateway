package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	postpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/post"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/SomeHowMicroservice/shm-be/gateway/request"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

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
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusCreated, "Tạo chủ đề bài viết thành công", gin.H{
		"topic_id": res.Id,
	})
}

func (h *PostHandler) GetAllTopicsAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.postClient.GetAllTopicsAdmin(ctx, &postpb.GetManyRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy danh sách chủ đề bài viết thành công", gin.H{
		"topics": res.Topics,
	})
}

func (h *PostHandler) GetAllTopics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.postClient.GetAllTopics(ctx, &postpb.GetManyRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy danh sách chủ đề bài viết thành công", gin.H{
		"topics": res.Topics,
	})
}

func (h *PostHandler) GetDeletedTopics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.postClient.GetDeletedTopics(ctx, &postpb.GetManyRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả chủ đề đã xóa thành công", gin.H{
		"topics": res.Topics,
	})
}

func (h *PostHandler) UpdateTopic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

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
	}); err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
	}

	common.JSON(c, http.StatusOK, "Cập nhật chủ đề bài viết thành công", nil)
}

func (h *PostHandler) DeleteTopic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	topicID := c.Param("id")

	if _, err := h.postClient.DeleteTopic(ctx, &postpb.DeleteOneRequest{
		Id:     topicID,
		UserId: user.Id,
	}); err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Chuyển chủ đề bài viết vào thùng rác thành công", nil)
}

func (h *PostHandler) DeleteTopics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	var req request.DeleteManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.postClient.DeleteTopics(ctx, &postpb.DeleteManyRequest{
		Ids:    req.IDs,
		UserId: user.Id,
	}); err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Chuyển danh sách chủ đề vào thùng rác thành công", nil)
}

func (h *PostHandler) RestoreTopic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	topicID := c.Param("id")

	if _, err := h.postClient.RestoreTopic(ctx, &postpb.RestoreOneRequest{
		Id:     topicID,
		UserId: user.Id,
	}); err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Khôi phục chủ đề bài viết thành công", nil)
}

func (h *PostHandler) RestoreTopics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	var req request.RestoreManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.postClient.RestoreTopics(ctx, &postpb.RestoreManyRequest{
		Ids:    req.IDs,
		UserId: user.Id,
	}); err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
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
	}); err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
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
	}); err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Xóa danh sách chủ đề thành công", nil)
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		common.JSON(c, http.StatusBadRequest, "Không thể parse form", nil)
		return
	}

	var req request.CreatePostForm

	req.Title = strings.TrimSpace(c.PostForm("title"))
	req.Content = strings.TrimSpace(c.PostForm("content"))
	req.TopicID = strings.TrimSpace(c.PostForm("post_id"))

	if isPublishedStr := c.PostForm("is_published"); isPublishedStr != "" {
		if isPublished, err := strconv.ParseBool(isPublishedStr); err == nil {
			req.IsPublished = &isPublished
		}
	}

	req.Images = []request.CreatePostImageForm{}
	imageCount := 0
	for {
		fileKey := fmt.Sprintf("images[%d][file]", imageCount)
		_, err := c.FormFile(fileKey)
		if err != nil {
			break
		}
		imageCount++
	}
	for j := 0; j < imageCount; j++ {
		isThumbnailKey := fmt.Sprintf("images[%d][is_thumbnail]", j)
		sortOrderKey := fmt.Sprintf("images[%d][sort_order]", j)
		fileKey := fmt.Sprintf("images[%d][file]", j)

		isThumbnailStr := strings.TrimSpace(c.PostForm(isThumbnailKey))
		sortOrderStr := strings.TrimSpace(c.PostForm(sortOrderKey))

		isThumbnail := false
		if isThumbnailStr != "" {
			isThumbnail, _ = strconv.ParseBool(isThumbnailStr)
		}

		sortOrder := 0
		if sortOrderStr != "" {
			sortOrder, _ = strconv.Atoi(sortOrderStr)
		}

		file, err := c.FormFile(fileKey)
		if err != nil {
			common.JSON(c, http.StatusBadRequest, fmt.Sprintf("Không tìm thấy file cho image %d: %s", j, err.Error()), nil)
			return
		}

		image := request.CreatePostImageForm{
			IsThumbnail: &isThumbnail,
			SortOrder:   sortOrder,
			File:        file,
		}

		req.Images = append(req.Images, image)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	images := make([]*postpb.CreateImageRequest, 0, len(req.Images))
	for _, img := range req.Images {
		openedFile, err := img.File.Open()
		if err != nil {
			common.JSON(c, http.StatusBadRequest, "Không mở được file", nil)
			return
		}
		defer openedFile.Close()

		fileBytes, err := io.ReadAll(openedFile)
		if err != nil {
			common.JSON(c, http.StatusInternalServerError, "Đọc file thất bại", nil)
			return
		}

		base64Data := base64.StdEncoding.EncodeToString(fileBytes)

		images = append(images, &postpb.CreateImageRequest{
			Base64Data:  base64Data,
			FileName:    img.File.Filename,
			IsThumbnail: *img.IsThumbnail,
			SortOrder:   int32(img.SortOrder),
		})
	}

	res, err := h.postClient.CreatePost(ctx, &postpb.CreatePostRequest{
		Title:       req.Title,
		Content:     req.Content,
		TopicId:     req.TopicID,
		IsPublished: *req.IsPublished,
		Images:      images,
		UserId:      user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusCreated, "Tạo bài viết thành công", gin.H{
		"post_id": res.Id,
	})
}
