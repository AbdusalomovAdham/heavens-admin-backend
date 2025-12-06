package departament

import "main/internal/entity"

type Create struct {
	Name   *entity.Name `json:"name"`
	Status *bool        `json:"status" default:"true"`
}

type Get struct {
	Id        int64  `json:"id"`
	Status    bool   `json:"status"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
}

type DepartamentById struct {
	Id        int64        `json:"id" bson:"id"`
	Status    bool         `json:"status" bson:"status"`
	CreatedAt string       `json:"created_at" bson:"created_at"`
	Name      *entity.Name `json:"name" bson:"name"`
}

type Update struct {
	Name   *entity.Name `json:"name"`
	Status *bool        `json:"status"`
}
