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

type DTOComments struct {
	ID      uuid.UUID
	FilmID  int
	Comment string
}

type DTOFavorites struct {
	ID     uuid.UUID
	FilmID int
	UserID uuid.UUID
}

type DTOAllFavorites struct {
	FilmID int
	Name   string
}
