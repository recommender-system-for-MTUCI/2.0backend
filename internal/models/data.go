package models

import "github.com/google/uuid"

type UserData struct {
	ID       uuid.UUID `json:"id"`
	IsAccess bool      `json:"is_access"`
}
