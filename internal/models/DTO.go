package models

import "github.com/google/uuid"

type DTORegister struct {
	ID           uuid.UUID
	Login        string
	Password     string
	Confirmation bool
	Code         int
}

type DTOPassword struct {
	Password string
}
