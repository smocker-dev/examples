package types

type Config struct {
	Port      int
	HotelsURL string
	UsersURL  string
	DSN       string
}

type Hotel struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Rooms int64  `json:"rooms"`
}

type Reservation struct {
	ID         int64 `bun:"id,pk,autoincrement" json:"id"`
	UserID     int64 `bun:"user_id,notnull" json:"user_id"`
	HotelID    int64 `bun:"hotel_id,notnull" json:"hotel_id"`
	RoomNumber int64 `bun:"room_number,notnull" json:"room_number"`
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CreateReservationBody struct {
	Hotel string `json:"hotel"`
	User  string `json:"user"`
	Rooms int64  `json:"rooms"`
}
