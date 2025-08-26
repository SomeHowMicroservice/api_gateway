package common

import "errors"

var (
	ErrUserNotFound = errors.New("không tìm thấy người dùng")

	ErrInvalidToken = errors.New("token không hợp lệ hoặc đã hết hạn")

	ErrRolesNotFound = errors.New("không tìm thấy roles")

	ErrUserIdNotFound = errors.New("không tìm thấy user_id")
)