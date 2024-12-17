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
	UserID  uuid.UUID
	Comment string
	Rating  float64
}

type DTOFavorites struct {
	ID     uuid.UUID
	FilmID int
	UserID uuid.UUID
}

type DTOAllFavorites struct {
	FilmID int
	Name   string
	Rating float64
}

type DTOFilmMain struct {
	FilmID int
	Name   string
	Rating float64
}
