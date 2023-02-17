package userclass

import (
	"fmt"
	"github.com/san035/basicApiGo/pkg/token"
	"time"
)

const (
	FormatDate = "02.01.2006 15:04:05"
	RoleUser   = "user"
	RoleAdmin  = "admin"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email" xml:"email" form:"email"`
	Password string `json:"password" xml:"password" form:"password"`
	Role     string `json:"role" xml:"role" form:"role"`
	CreatAt  int64  `json:"CreatAt"`  // Дата создания UNIX
	UpdateAt int64  `json:"UpdateAt"` // Дата обновления UNIX
	Exp      int64  `json:"exp"`      // Дата экспирации
}

// MarshalJSON переопределение CreatAt и UpdateAt для вызовов Marshal
func (user *User) MarshalJSON() ([]byte, error) {
	userStr := fmt.Sprintf("{\"id\":\"%s\",\"email\":\"%s\",\"role\":\"%s\",\"CreatAt\":\"%s\",\"UpdateAt\":\"%s\"}",
		user.ID, user.Email, user.Role, time.Unix(user.CreatAt, 0).Format(FormatDate), time.Unix(user.UpdateAt, 0).Format(FormatDate))
	return []byte(userStr), nil
}

// Создание токена
func (user *User) CreateToken() (string, error) {
	return token.Create(map[string]interface{}{"email": user.Email, "role": user.Role, "id": user.ID})
}
