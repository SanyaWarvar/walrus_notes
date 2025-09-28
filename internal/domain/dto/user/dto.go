package user

import (
	"wn/internal/infrastructure/repository/user"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	ImgUrl    string    `json:"imgUrl"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

func UserDtoFromEntity(entity *user.User) *User {
	return &User{
		Id:        entity.Id,
		Username:  entity.Username,
		Email:     entity.Email,
		ImgUrl:    entity.ImgUrl,
		Role:      entity.Role,
		CreatedAt: entity.CreatedAt,
	}
}
