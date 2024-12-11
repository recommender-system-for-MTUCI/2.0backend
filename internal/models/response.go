package models

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
type Comments struct {
	UserLogin string `json:"user_login"`
	Com       string `json:"comment"`
}
type FavouriteResponse struct {
	Name string `json:"name"`
}

type ResponseFilm struct {
	Title               string     `json:"title"`
	Genres              []string   `json:"genres"`
	Overview            string     `json:"overview"`
	ProductionCompanies []string   `json:"production_companies"`
	ProductionCountries []string   `json:"production_countries"`
	ReleaseDate         string     `json:"release_date"`
	RunTime             int        `json:"run_time"`
	VoteAverage         float64    `json:"vote_average"`
	VoteCount           int        `json:"vote_count"`
	Actor               []string   `json:"actor"`
	KeyWords            []string   `json:"key_words"`
	Director            string     `json:"director"`
	WeightRating        float64    `json:"weight_rating"`
	FilmsComments       []Comments `json:"comments"`
}
type ResponseAllFavorites struct {
	FilmID int     `json:"film_id"`
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
}
