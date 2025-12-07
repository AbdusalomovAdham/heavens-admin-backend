package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:roles"`

	Id        int64      `json:"id"`
	Name      *Name      `json:"name" bun:"name,type:jsonb"`
	CreatedBy int64      `json:"created_by" bun:"created_by,default:null"`
	Status    bool       `json:"status" bun:"status,default:true"`
	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by,default:null"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at,default:null"`
	DeletedBy *int64     `json:"deleted_by" bun:"deleted_by,default:null"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at,default:null"`
}
