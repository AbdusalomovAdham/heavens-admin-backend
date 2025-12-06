package room

type Create struct {
	Corpus     uint   `json:"corpus"`
	EmployeeId int64  `json:"employee_id"`
	RoomNumber uint64 `json:"room_number"`
	RoomType   int64  `json:"room_type"`
}

type RoomPreview struct {
	Id                 int64   `json:"id" bun:"id"`
	EmployeeId         int64   `json:"employee_id" bun:"employee_id"`
	EmployeeFirstName  string  `json:"employee_first_name" bun:"first_name"`
	EmployeeLastName   string  `json:"employee_last_name" bun:"last_name"`
	EmployeeMiddleName string  `json:"employee_middle_name" bun:"middle_name"`
	RoomNumber         uint64  `json:"room_number" bun:"room_number"`
	RoomType           int64   `json:"room_type" bun:"room_type"`
	QRPath             *string `json:"qr_path" bun:"path,default:null"`
	Corpus             uint    `json:"corpus" bun:"corpus"`
	Status             bool    `json:"status" bun:"status"`
	CreatedAt          string  `json:"created_at" bun:"created_at"`
}
