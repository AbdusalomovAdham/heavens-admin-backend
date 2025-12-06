package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Room struct {
	bun.BaseModel `bun:"table:rooms"`

	Id         int64      `json:"id" bun:"id,pk,autoincrement"`
	EmployeeId int64      `json:"employee_id" bun:"employee_id,default:null"`
	RoomNumber uint64     `json:"room_number" bun:"room_number"`
	RoomType   uint64     `json:"room_type" bun:"room_type"`
	Corpus     uint       `json:"corpus" bun:"corpus"`
	Status     bool       `json:"status" bun:"status,default:true"`
	CreatedBy  int64      `json:"created_by" bun:"created_by,default:null"`
	CreatedAt  time.Time  `json:"created_at" bun:"created_at"`
	UpdatedBy  *int64     `json:"updated_by" bun:"updated_by,default:null"`
	UpdatedAt  *time.Time `json:"updated_at" bun:"updated_at,default:null"`
	DeletedBy  *int64     `json:"deleted_by" bun:"deleted_by,default:null"`
	DeletedAt  *time.Time `json:"deleted_at" bun:"deleted_at,default:null"`
}

//Role

//1-super amdin
//2-admin
//3-user
