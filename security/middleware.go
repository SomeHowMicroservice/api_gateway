package security

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/SomeHowMicroservice/gateway/common"
	userpb "github.com/SomeHowMicroservice/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequireRefreshToken(refreshName, secretKey string, userClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(refreshName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: common.ErrUnAuth.Error(),
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

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		userRes, err := fetchUserFromUserService(ctx, userID, userClient)
		if err != nil {
			switch err {
			case common.ErrUserNotFound:
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
					Message: err.Error(),
				})
				return
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{
					Message: err.Error(),
				})
				return
			}
		}

		if !slices.Equal(userRes.Roles, userRoles) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: common.ErrInvalidUser.Error(),
			})
			return
		}

		c.Set("user_id", userRes.Id)
		c.Set("user_roles", userRes.Roles)
		c.Next()
	}
}

func RequireAuth(accessName, secretKey string, userClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(accessName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: common.ErrUnAuth.Error(),
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

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		userRes, err := fetchUserFromUserService(ctx, userID, userClient)
		if err != nil {
			switch err {
			case common.ErrUserNotFound:
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
					Message: err.Error(),
				})
				return
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{
					Message: err.Error(),
				})
				return
			}
		}

		if !slices.Equal(userRes.Roles, userRoles) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: common.ErrInvalidUser.Error(),
			})
			return
		}

		if !hasRoleUser(common.RoleUser, userRes.Roles) {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ApiResponse{
				Message: common.ErrForbidden.Error(),
			})
			return
		}

		c.Set("user", userRes)
		c.Next()
	}
}

func OptionalAuth(accessName, secretKey string, userClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(accessName)
		if err != nil || tokenStr == "" {
			c.Set("user_id", "")
			c.Next()
			return
		}

		claims, err := ParseToken(tokenStr, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: err.Error(),
			})
			return
		}

		userID, _, err := ExtractToken(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		res, err := userClient.CheckUserExistsById(ctx, &userpb.CheckUserExistsByIdRequest{
			Id: userID,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{
					Message: st.Message(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ApiResponse{
				Message: err.Error(),
			})
			return
		}

		if !res.Exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: common.ErrUnAuth.Error(),
			})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func RequireMultiRoles(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: common.ErrUnAuth.Error(),
			})
			return
		}

		user := userAny.(*userpb.UserPublicResponse)

		if !hasAtLeastOneRole(user.Roles, allowedRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ApiResponse{
				Message: common.ErrForbidden.Error(),
			})
			return
		}

		c.Next()
	}
}

func RequireSingleRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: common.ErrUnAuth.Error(),
			})
			return
		}

		user := userAny.(*userpb.UserPublicResponse)

		if len(user.Roles) > 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ApiResponse{
				Message: common.ErrForbidden.Error(),
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
				return nil, fmt.Errorf("%s", st.Message())
			}
		}
		return nil, err
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
