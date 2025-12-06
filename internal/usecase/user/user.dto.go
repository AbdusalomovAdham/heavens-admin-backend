package user

import (
	"time"
)

type Create struct {
	Avatar         string
	FirstName      string     `json:"first_name" form:"first_name"`
	LastName       string     `json:"last_name" form:"last_name"`
	MiddleName     string     `json:"middle_name" form:"middle_name"`
	RegionId       *uint16    `json:"region_id" form:"region_id"`
	DistrictId     *uint16    `json:"district_id" form:"district_id"`
	GenderId       *uint8     `json:"gender_id" form:"gender_id"`
	BirthDate      *time.Time `json:"birth_date" form:"birth_date"`
	Login          string     `json:"login" form:"login"`
	Password       string     `json:"password" form:"password"`
	Email          *string    `json:"email" form:"email"`
	PhoneNumber    *string    `json:"phone_number" form:"phone_number"`
	MobileNumber   *string    `json:"mobile_number" form:"mobile_number"`
	Role           *uint      `json:"role" form:"role"`
	ManagementId   *uint      `json:"management_id" form:"management_id"`
	PositionId     *uint      `json:"position_id" form:"position_id"`
	WorkStatusId   *uint      `json:"work_status_id" form:"work_status_id"`
	Salary         *uint32    `json:"salary" form:"salary"`
	PassportNumber *uint64    `json:"passport_number" form:"passport_number"`
	JSHSHIR        *uint64    `json:"jshshir" form:"jshshir"`
	PassportScan   *string    `json:"passport_scan"`
	CarPrefix      *int8      `json:"car_prefix" form:"car_prefix"`
	CarNumber      *string    `json:"car_number" form:"car_number"`
	DimlomaFile    *string    `json:"diploma_file"`
	CVFile         *string    `json:"cv_file"`
	Status         *bool      `json:"status" form:"status"`
	CreatedBy      *int64     `json:"created_by" form:"created_by"`
}

type Update struct {
	Avatar         *string
	FirstName      *string    `json:"first_name" form:"first_name"`
	LastName       *string    `json:"last_name" form:"last_name"`
	MiddleName     *string    `json:"middle_name" form:"middle_name"`
	RegionId       *uint16    `json:"region_id" form:"region_id"`
	DistrictId     *uint16    `json:"district_id" form:"district_id"`
	GenderId       *uint8     `json:"gender_id" form:"gender_id"`
	BirthDate      *time.Time `json:"birth_date" form:"birth_date"`
	Login          *string    `json:"login" form:"login"`
	Password       *string    `json:"password" form:"password"`
	Email          *string    `json:"email" form:"email"`
	PhoneNumber    *string    `json:"phone_number" form:"phone_number"`
	MobileNumber   *string    `json:"mobile_number" form:"mobile_number"`
	Role           *uint      `json:"role" form:"role"`
	ManagementId   *uint      `json:"management_id" form:"management_id"`
	PositionId     *uint      `json:"position_id" form:"position_id"`
	WorkStatusId   *uint      `json:"work_status_id" form:"work_status_id"`
	Salary         *uint32    `json:"salary" form:"salary"`
	PassportNumber *uint64    `json:"passport_number" form:"passport_number"`
	JSHSHIR        *uint64    `json:"jshshir" form:"jshshir"`
	PassportScan   *string    `json:"passport_scan"`
	CarPrefix      *int8      `json:"car_prefix" form:"car_prefix"`
	CarNumber      *string    `json:"car_number" form:"car_number"`
	DimlomaFile    *string    `json:"diploma_file"`
	CVFile         *string    `json:"cv_file"`
	Status         *bool      `json:"status" form:"status"`
	UpdatedBy      *int64     `json:"created_by" form:"created_by"`
}

type Filter struct {
	Limit  *int
	Offset *int
	Role   *int
	Order  *string
	Status *bool
}

type UserPreview struct {
	Id           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Avatar       string `json:"avatar"`
	ManagementId int64  `json:"management_id"`
	PositionId   int64  `json:"position_id"`
	Status       bool   `json:"status"`
}
