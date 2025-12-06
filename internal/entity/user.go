package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	Id             int64      `json:"id" bun:"id,pk,autoincrement"`
	Avatar         string     `json:"avatar" bun:"avatar"`
	FirstName      string     `json:"first_name" bun:"first_name,notnull"`
	LastName       string     `json:"last_name" bun:"last_name,notnull"`
	MiddleName     string     `json:"middle_name" bun:"middle_name,notnull"`
	RegionId       *uint16    `json:"region_id" bun:"region_id,notnull"`
	DistrictId     *uint16    `json:"district_id" bun:"district_id,notnull"`
	GenderId       *uint8     `json:"gender_id" bun:"gender_id,notnull"`
	BirthDate      *time.Time `json:"birth_date" bun:"birth_date,notnull"`
	Login          string     `json:"login" bun:"login,notnull"`
	Password       string     `json:"password" bun:"password,notnull"`
	Email          *string    `json:"email" bun:"email,notnull"`
	PhoneNumber    *string    `json:"phone_number" bun:"phone_number,notnull"`
	MobileNumber   *string    `json:"mobile_number" bun:"mobile_number,notnull"`
	Role           *uint      `json:"role" bun:"role,default:3"`
	ManagementId   *uint      `json:"management_id" bun:"management_id,notnull"`
	PositionId     *uint      `json:"position_id" bun:"position_id,notnull"`
	WorkStatusId   *uint      `json:"work_status_id" bun:"work_status_id,notnull"`
	Salary         *uint32    `json:"salary" bun:"salary,notnull"`
	PassportNumber *uint64    `json:"passport_number" bun:"passport_number,notnull"`
	JSHSHIR        *uint64    `json:"jshshir" bun:"jshshir,notnull"`
	PassportScan   *string    `json:"passport_scan" bun:"passport_scan,notnull"`
	CarPrefix      *int8      `json:"car_prefix" bun:"car_prefix,notnull"`
	CarNumber      *string    `json:"car_number" bun:"car_number,notnull"`
	DimlomaFile    *string    `json:"diploma_file" bun:"diploma_file,notnull"`
	CVFile         *string    `json:"cv_file" bun:"cv_file,notnull"`
	Status         bool       `json:"status" bun:"status,default:true"`
	CreatedBy      int64      `json:"created_by" bun:"created_by,default:null"`
	CreatedAt      time.Time  `json:"created_at" bun:"created_at"`
	UpdatedBy      int64      `json:"updated_by" bun:"updated_by,default:null"`
	UpdatedAt      time.Time  `json:"updated_at" bun:"updated_at,default:null"`
	DeletedBy      int64      `json:"deleted_by" bun:"deleted_by,default:null"`
	DeletedAt      time.Time  `json:"deleted_at" bun:"deleted_at,default:null"`
}

//Role

//1-super amdin
//2-admin
//3-user
