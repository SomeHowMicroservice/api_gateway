package request

import "time"

type UpdateProfileRequest struct {
	FirstName *string    `json:"first_name" binding:"omitempty"`
	LastName  *string    `json:"last_name" binding:"omitempty"`
	Gender    *string    `json:"gender" binding:"omitempty,oneof=male female"`
	DOB       *time.Time `json:"dob" binding:"omitempty"`
}

type UpdateMeasurementRequest struct {
	Height *int `json:"height" binding:"omitempty"`
	Weight *int `json:"weight" binding:"omitempty"`
	Chest  *int `json:"chest" binding:"omitempty"`
	Waist  *int `json:"waist" binding:"omitempty"`
	Butt   *int `json:"butt" binding:"omitempty"`
}

type CreateMyAddressRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Street      string `json:"street" binding:"required"`
	Ward        string `json:"ward" binding:"required"`
	Province    string `json:"province" binding:"required"`
	IsDefault   *bool   `json:"is_default" binding:"required"`
}

type UpdateAddressRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Street      string `json:"street" binding:"required"`
	Ward        string `json:"ward" binding:"required"`
	Province    string `json:"province" binding:"required"`
	IsDefault   *bool   `json:"is_default" binding:"required"`
}
