package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `db:"id"`
	Username  string     `db:"username"`
	Password  string     `db:"password"`
	DeletedAt *time.Time `db:"deleted_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	CreatedAt time.Time  `db:"created_at"`
}

func (u User) TableName() string {
	return "users"
}

func (u User) Fields() string {
	fields := []string{"id", "username", "password", "deleted_at", "updated_at", "created_at"}
	for ind, field := range fields {
		fields[ind] = fmt.Sprintf("%s.%s", u.TableName(), field)
	}

	return strings.Join(fields, ", ")
}
