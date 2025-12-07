package entity

type Filter struct {
	Limit  *int
	Offset *int
	Role   *int
	Order  *string
	Status *bool
}
