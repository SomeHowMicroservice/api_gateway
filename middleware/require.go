package middleware

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequireRefreshToken(refreshName string, secretKey string, userClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(refreshName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: "Không tìm thấy token làm mới",
			})
			return
		}

		claims, err := ParseToken(tokenStr, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: err.Error(),
			})
			return
		}

		userID, userRoles, err := ExtractToken(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: err.Error(),
			})
			return
		}

		ctx := c.Request.Context()
		userRes, err := fetchUserFromUserService(ctx, userID, userClient)
		if err != nil {
			switch err {
			case common.ErrUserNotFound:
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
					Message: err.Error(),
				})
				return
			default:
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
					Message: err.Error(),
				})
				return
			}
		}

		if !slices.Equal(userRes.Roles, userRoles) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: "người dùng không hợp lệ",
			})
			return
		}

		c.Set("user_id", userRes.Id)
		c.Set("user_roles", userRes.Roles)
		c.Next()
	}
}

func RequireAuth(accessName string, secretKey string, userClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(accessName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: "Không tìm thấy token",
			})
			return
		}

		claims, err := ParseToken(tokenStr, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: err.Error(),
			})
			return
		}

		userID, userRoles, err := ExtractToken(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: err.Error(),
			})
			return
		}

		ctx := c.Request.Context()
		userRes, err := fetchUserFromUserService(ctx, userID, userClient)
		if err != nil {
			switch err {
			case common.ErrUserNotFound:
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
					Message: err.Error(),
				})
				return
			default:
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
					Message: err.Error(),
				})
				return
			}
		}

		if !slices.Equal(userRes.Roles, userRoles) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: "người dùng không hợp lệ",
			})
			return
		}

		if !hasRoleUser(common.RoleUser, userRes.Roles) {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ApiResponse{
				Message: "không có quyền truy cập",
			})
			return
		}

		c.Set(common.RoleUser, userRes)
		c.Next()
	}
}

func RequireMultiRoles(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get(common.RoleUser)
		if !exists {
			common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
			return
		}

		user, ok := userAny.(*userpb.UserPublicResponse)
		if !ok {
			common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
			return
		}

		if !hasAtLeastOneRole(user.Roles, allowedRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ApiResponse{
				Message: "không có quyền truy cập",
			})
			return
		}

		c.Next()
	}
}

func hasRoleUser(role string, roles []string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

func fetchUserFromUserService(ctx context.Context, userID string, userClient userpb.UserServiceClient) (*userpb.UserPublicResponse, error) {
	userRes, err := userClient.GetUserPublicById(ctx, &userpb.GetOneRequest{Id: userID})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, common.ErrUserNotFound
			default:
				return nil, fmt.Errorf("lỗi user service: %s", st.Message())
			}
		}
		return nil, fmt.Errorf("lỗi hệ thống: %w", err)
	}

	return userRes, nil
}

func hasAtLeastOneRole(userRoles, allowedRoles []string) bool {
	roleSet := make(map[string]struct{})
	for _, r := range userRoles {
		roleSet[r] = struct{}{}
	}

	for _, allowed := range allowedRoles {
		if _, ok := roleSet[allowed]; ok {
			return true
		}
	}

	return false
}
