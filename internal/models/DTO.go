package models

import "github.com/google/uuid"

type DTOLogin struct {
	ID       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
}