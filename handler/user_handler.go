package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/request"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userClient userpb.UserServiceClient
}

func NewUserHandler(userClient userpb.UserServiceClient) *UserHandler {
	return &UserHandler{userClient}
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
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

	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	profileID := c.Param("id")
	if profileID != user.Profile.Id {
		common.JSON(c, http.StatusForbidden, "Không có quyền chỉnh sửa", nil)
		return
	}

	var firstName, lastName, gender, dob string
	if req.FirstName != nil {
		firstName = *req.FirstName
	}
	if req.LastName != nil {
		lastName = *req.LastName
	}
	if req.Gender != nil {
		gender = *req.Gender
	}
	if req.DOB != nil {
		dob = req.DOB.Format("2006-01-02")
	}

	res, err := h.userClient.UpdateProfile(ctx, &userpb.UpdateProfileRequest{
		Id:        profileID,
		FirstName: firstName,
		LastName:  lastName,
		Gender:    gender,
		Dob:       dob,
		UserId:    user.Id,
	})
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

	common.JSON(c, http.StatusOK, "Cập nhật hồ sơ thành công", gin.H{
		"user": ToAuthResponse(res),
	})
}

func (h *UserHandler) MyMeasurements(c *gin.Context) {
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

	res, err := h.userClient.GetMeasurementByUserId(ctx, &userpb.GetByUserIdRequest{
		UserId: user.Id,
	})
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

	common.JSON(c, http.StatusOK, "Lấy độ đo người dùng thành công", gin.H{
		"measurement": res,
	})
}

func (h *UserHandler) UpdateMeasurement(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.UpdateMeasurementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

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

	measurementID := c.Param("id")
	var height, weight, chest, waist, butt int32
	if req.Height != nil {
		height = int32(*req.Height)
	}
	if req.Weight != nil {
		weight = int32(*req.Weight)
	}
	if req.Chest != nil {
		chest = int32(*req.Chest)
	}
	if req.Waist != nil {
		waist = int32(*req.Waist)
	}
	if req.Butt != nil {
		butt = int32(*req.Butt)
	}

	res, err := h.userClient.UpdateMeasurement(ctx, &userpb.UpdateMeasurementRequest{
		Id:     measurementID,
		Height: height,
		Weight: weight,
		Chest:  chest,
		Waist:  waist,
		Butt:   butt,
		UserId: user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.PermissionDenied:
				common.JSON(c, http.StatusForbidden, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Cập nhật độ đo thành công", gin.H{
		"measurement": res,
	})
}

func (h *UserHandler) MyAddresses(c *gin.Context) {
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

	res, err := h.userClient.GetAddressesByUserId(ctx, &userpb.GetByUserIdRequest{
		UserId: user.Id,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy địa chỉ người dùng thành công", res)
}

func (h *UserHandler) CreateMyAddress(c *gin.Context) {
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

	var req request.CreateMyAddressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	var isDefault bool
	if req.IsDefault != nil {
		isDefault = *req.IsDefault
	}

	res, err := h.userClient.CreateAddress(ctx, &userpb.CreateAddressRequest{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Street:      req.Street,
		Ward:        req.Ward,
		Province:    req.Province,
		IsDefault:   isDefault,
		UserId:      user.Id,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Tạo địa chỉ thành công", gin.H{
		"address": res,
	})
}

func (h *UserHandler) UpdateAddress(c *gin.Context) {
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

	addressID := c.Param("id")
	var req request.UpdateAddressRequest
	if err := c.ShouldBind(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	var isDefault bool
	if req.IsDefault != nil {
		isDefault = *req.IsDefault
	}

	res, err := h.userClient.UpdateAddress(ctx, &userpb.UpdateAddressRequest{
		Id:          addressID,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Street:      req.Street,
		Ward:        req.Ward,
		Province:    req.Province,
		IsDefault:   isDefault,
		UserId:      user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.PermissionDenied:
				common.JSON(c, http.StatusForbidden, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Cập nhật địa chỉ thành công", gin.H{
		"address": res,
	})
}

func (h *UserHandler) DeleteAddress(c *gin.Context) {
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

	addressID := c.Param("id")

	_, err := h.userClient.DeleteAddress(ctx, &userpb.DeleteAddressRequest{
		Id:     addressID,
		UserId: user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.PermissionDenied:
				common.JSON(c, http.StatusForbidden, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Xoá địa chỉ thành công", nil)
}
