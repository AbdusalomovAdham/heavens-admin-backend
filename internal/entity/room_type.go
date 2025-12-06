package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type RoomType struct {
	bun.BaseModel `bun:"table:room_types"`

	Id        int64      `json:"id" bun:"id,pk,autoincrement"`
	Name      string     `json:"name" bun:"name,default:null"`
	Status    bool       `json:"status" bun:"status,default:true"`
	CreatedBy int64      `json:"created_by" bun:"created_by,default:null"`
	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by,default:null"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at,default:null"`
}
