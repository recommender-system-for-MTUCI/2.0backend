package models

import "github.com/google/uuid"

type ResponseRegister struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type ResponseLogin struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Comment struct {
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
}
type ResponseComments struct {
	ID        uuid.UUID `json:"id"`
	Rating    float64   `json:"rating"`
	UserLogin string    `json:"user_login"`
	Comment   string    `json:"comment"`
}
type FavouriteResponse struct {
	Name string `json:"name"`
}

type ResponseFilm struct {
	Title               string             `json:"title"`
	Genres              []string           `json:"genres"`
	Overview            string             `json:"overview"`
	ProductionCompanies []string           `json:"production_companies"`
	ProductionCountries []string           `json:"production_countries"`
	ReleaseDate         string             `json:"release_date"`
	RunTime             int                `json:"run_time"`
	VoteAverage         float64            `json:"vote_average"`
	VoteCount           int                `json:"vote_count"`
	Actor               []string           `json:"actor"`
	KeyWords            []string           `json:"key_words"`
	Director            string             `json:"director"`
	WeightRating        float64            `json:"weight_rating"`
	FilmsComments       []ResponseComments `json:"comments"`
}
type ResponseAllFavorites struct {
	FilmID int     `json:"film_id"`
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
}

type ResponseFilmMain struct {
	FilmID int     `json:"film_id"`
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
}
type ResponseFilmGenre struct {
	FilmID int     `json:"film_id"`
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
}

type ResponseID struct {
	RecommendMovies []int `json:"recommend_movies"`
}

type ResponseFilmPage struct {
	Film  ResponseFilm       `json:"film"`
	Films []ResponseFilmName `json:"films"`
}
type ResponseFilmName struct {
	FilmID int     `json:"film_id"`
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
}
