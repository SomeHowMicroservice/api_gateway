package common

import "errors"

var (
	ErrUserNotFound = errors.New("không tìm thấy người dùng")

	ErrInvalidToken = errors.New("token không hợp lệ hoặc đã hết hạn")

	ErrInvalidUser = errors.New("người dùng không hợp lệ")

	ErrRolesNotFound = errors.New("không tìm thấy các quyền")

	ErrUserIdNotFound = errors.New("không tìm thấy user_id")

	ErrForbidden = errors.New("không có quyền truy cập")

	ErrUnAuth = errors.New("bạn chưa đăng nhập")
)