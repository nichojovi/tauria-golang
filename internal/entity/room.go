package entity

type Room struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	HostUser     string `json:"host_user"`
	Participants string `json:"participant"`
	Capacity     int64  `json:"capacity"`
}

type RoomDB struct {
	ID           int64  `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	HostUser     string `json:"host_user" db:"host_user"`
	Participants []byte `json:"participant" db:"participant"`
	Capacity     int64  `json:"capacity" db:"capacity"`
}
