package roomtype

type Create struct {
	Name   string `json:"name"`
	Status *bool  `json:"status"`
}
