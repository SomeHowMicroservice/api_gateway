package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/request"
	authpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/auth"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authClient authpb.AuthServiceClient
	cfg        *config.AppConfig
}

func NewAuthHandler(authClient authpb.AuthServiceClient, cfg *config.AppConfig) *AuthHandler {
	return &AuthHandler{
		authClient,
		cfg,
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.authClient.SignUp(ctx, &authpb.SignUpRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			case codes.InvalidArgument:
				common.JSON(c, http.StatusBadRequest, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Đăng ký thành công, vui lòng kiểm tra email để lấy mã xác thực", gin.H{
		"registration_token": res.RegistrationToken,
	})
}

func (h *AuthHandler) VerifySignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.VerifySignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.authClient.VerifySignUp(ctx, &authpb.VerifySignUpRequest{
		RegistrationToken: req.RegistrationToken,
		Otp:               req.Otp,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.InvalidArgument:
				common.JSON(c, http.StatusBadRequest, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.SetCookie(h.cfg.Jwt.AccessName, res.AccessToken, int(res.AccessExpiresIn), "/", "", false, true)
	c.SetCookie(h.cfg.Jwt.RefreshName, res.RefreshToken, int(res.RefreshExpiresIn), "/api/v1/auth/refresh", "", false, true)

	common.JSON(c, http.StatusCreated, "Đăng ký thành công", gin.H{
		"user": res.User,
	})
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.authClient.SignIn(ctx, &authpb.SignInRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.InvalidArgument:
				common.JSON(c, http.StatusBadRequest, st.Message(), nil)
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

	c.SetCookie(h.cfg.Jwt.AccessName, res.AccessToken, int(res.AccessExpiresIn), "/", "", false, true)
	c.SetCookie(h.cfg.Jwt.RefreshName, res.RefreshToken, int(res.RefreshExpiresIn), "/api/v1/auth/refresh", "", false, true)

	common.JSON(c, http.StatusOK, "Đăng nhập thành công", gin.H{
		"user": res.User,
	})
}

func (h *AuthHandler) SignOut(c *gin.Context) {
	c.SetCookie(h.cfg.Jwt.AccessName, "", 0, "/", "", false, true)
	c.SetCookie(h.cfg.Jwt.RefreshName, "", 0, "/api/v1/auth/refresh", "", false, true)

	common.JSON(c, http.StatusOK, "Đăng xuất thành công", nil)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
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

	common.JSON(c, http.StatusOK, "Lấy thông tin người dùng thành công", gin.H{
		"user": ToAuthResponse(user),
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	// Lấy UserID từ Context ra
	userIDAny, exists := c.Get("user_id")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có id người dùng", nil)
		return
	}
	userID, ok := userIDAny.(string)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi id người dùng", nil)
		return
	}
	// Lấy User Roles từ Context ra
	userRolesAny, exists := c.Get("user_roles")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có các quyền người dùng", nil)
		return
	}
	userRoles, ok := userRolesAny.([]string)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển các quyền người dùng", nil)
		return
	}
	res, err := h.authClient.RefreshToken(ctx, &authpb.RefreshTokenRequest{
		Id:    userID,
		Roles: userRoles,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.SetCookie(h.cfg.Jwt.AccessName, res.AccessToken, int(res.AccessExpiresIn), "/", "", false, true)
	c.SetCookie(h.cfg.Jwt.RefreshName, res.RefreshToken, int(res.RefreshExpiresIn), "/api/v1/auth/refresh", "", false, true)

	common.JSON(c, http.StatusOK, "Làm mới token thành công", nil)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
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
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}
	res, err := h.authClient.ChangePassword(ctx, &authpb.ChangePasswordRequest{
		Id:          user.Id,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.InvalidArgument:
				common.JSON(c, http.StatusBadRequest, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.SetCookie(h.cfg.Jwt.AccessName, res.AccessToken, int(res.AccessExpiresIn), "/", "", false, true)
	c.SetCookie(h.cfg.Jwt.RefreshName, res.RefreshToken, int(res.RefreshExpiresIn), "/api/v1/auth/refresh", "", false, true)

	common.JSON(c, http.StatusOK, "Đổi mật khẩu thành thành công", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.authClient.ForgotPassword(ctx, &authpb.ForgotPasswordRequest{
		Email: req.Email,
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

	common.JSON(c, http.StatusOK, "Xác thực quên mật khẩu thành công, vui lòng kiểm tra email để lấy mã xác thực", gin.H{
		"forgot_password_token": res.ForgotPasswordToken,
	})
}

func (h *AuthHandler) VerifyForgotPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.VerifyForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.authClient.VerifyForgotPassword(ctx, &authpb.VerifyForgotPasswordRequest{
		ForgotPasswordToken: req.ForgotPasswordToken,
		Otp:                 req.Otp,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.InvalidArgument:
				common.JSON(c, http.StatusBadRequest, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Xác thực quên mật khẩu thành công", gin.H{
		"reset_password_token": res.ResetPasswordToken,
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	_, err := h.authClient.ResetPassword(ctx, &authpb.ResetPasswordRequest{
		ResetPasswordToken: req.ResetPasswordToken,
		NewPassword:        req.NewPassword,
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

	common.JSON(c, http.StatusOK, "Làm mới mật khẩu thành công, vui lòng đăng nhập lại", nil)
}

func (h *AuthHandler) AdminSignIn(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.authClient.AdminSignIn(ctx, &authpb.SignInRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.InvalidArgument:
				common.JSON(c, http.StatusBadRequest, st.Message(), nil)
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

	c.SetCookie(h.cfg.Jwt.AccessName, res.AccessToken, int(res.AccessExpiresIn), "/", "", false, true)
	c.SetCookie(h.cfg.Jwt.RefreshName, res.RefreshToken, int(res.RefreshExpiresIn), "/api/v1/auth/refresh", "", false, true)

	common.JSON(c, http.StatusOK, "Đăng nhập thành công", gin.H{
		"user": res.User,
	})
}

func ToAuthResponse(userRes *userpb.UserPublicResponse) *authpb.AuthResponse {
	return &authpb.AuthResponse{
		Id:        userRes.Id,
		Username:  userRes.Username,
		Email:     userRes.Email,
		CreatedAt: userRes.CreatedAt,
		Profile: &authpb.ProfileResponse{
			Id:        userRes.Profile.Id,
			FirstName: userRes.Profile.FirstName,
			LastName:  userRes.Profile.LastName,
			Gender:    userRes.Profile.Gender,
			Dob:       userRes.Profile.Dob,
		},
	}
}
