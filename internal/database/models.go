// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      sql.NullString
}
