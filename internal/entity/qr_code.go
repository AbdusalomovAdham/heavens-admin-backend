package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type QRCode struct {
	bun.BaseModel `bun:"table:qr_codes"`

	Id        int64      `json:"id" bun:"id,pk,autoincrement"`
	RoomId    int64      `json:"room_id" bun:"room_id,default:null"`
	Path      string     `json:"path" bun:"path,default:null"`
	Status    bool       `json:"status" bun:"status,default:true"`
	CreatedBy int64      `json:"created_by" bun:"created_by,default:null"`
	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	UpdatedBy *int64     `json:"updated_by" bun:"updated_by,default:null"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at,default:null"`
}
