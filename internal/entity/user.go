package entity

var (
	MinCharacter = 5
)

type User struct {
	UserName string `json:"user_name" db:"user_name"`
	Password string `json:"password" db:"password"`
	FullName string `json:"full_name" db:"full_name"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`
}
