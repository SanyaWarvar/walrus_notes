package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	Role           string    `json:"role"`
	ImgUrl         string    `json:"imgUrl"`
	ConfirmedEmail bool      `json:"confirmedEmail"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Role struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
