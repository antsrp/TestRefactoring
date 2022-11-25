package user

import "time"

type (
	User struct {
		CreatedAt   time.Time `json:"created_at"`
		DisplayName string    `json:"display_name"`
		Email       string    `json:"email"`
	}
	UserList  map[int]User
	UserStore struct {
		Increment int      `json:"increment"`
		List      UserList `json:"list"`
	}
)

func CreateUser(displayName, email string) *User {
	return &User{
		CreatedAt:   time.Now(),
		DisplayName: displayName,
		Email:       email,
	}
}

type Storage interface {
	Add(*User) int
	Delete(int) error
	Update(int, string) error
	GetUser(int) (*User, error)
	GetUsers() *UserStore
	Save() error
}
