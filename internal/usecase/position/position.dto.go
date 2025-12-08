package position

import "main/internal/entity"

type Create struct {
	Name         *entity.Name `json:"name"`
	Status       *bool        `json:"status" default:"true"`
	DepartmentId int64        `json:"department_id"`
}

type Get struct {
	Id           int64  `json:"id"`
	Status       bool   `json:"status"`
	CreatedAt    string `json:"created_at"`
	DepartmentId int64  `json:"department_id" bson:"department_id"`
	Name         string `json:"name"`
}

type PositionById struct {
	Id           int64        `json:"id" bson:"id"`
	Status       bool         `json:"status" bson:"status"`
	CreatedAt    string       `json:"created_at" bson:"created_at"`
	DepartmentId int64        `json:"department_id" bson:"department_id"`
	Name         *entity.Name `json:"name" bson:"name"`
}

type Update struct {
	Name   *entity.Name `json:"name"`
	Status *bool        `json:"status"`
}
