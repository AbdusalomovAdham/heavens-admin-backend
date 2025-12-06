package qrcode

type Create struct {
	RoomId int64  `json:"room_id"`
	Url    string `json:"url"`
}
