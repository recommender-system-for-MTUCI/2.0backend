package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
)

type UserStorage interface {
	AddUserInDB(ctx context.Context, user *models.DTORegister) error
	GetCodeFromDB(ctx context.Context, id uuid.UUID) (int, error)
	GetPasswordFromDB(ctx context.Context, id uuid.UUID) (string, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error
	GetStatusFromUser(ctx context.Context, id uuid.UUID) (bool, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUserIdByEmail(ctx context.Context, email string) (uuid.UUID, string, error)
	UpdateUserStatus(ctx context.Context, id uuid.UUID) error
	GetMe(ctx context.Context, id uuid.UUID) (string, error)
}
type FilmStorage interface {
	AddNewComment(ctx context.Context, data *models.DTOComments) error
	AddFilmToFavourites(ctx context.Context, data *models.DTOFavorites) error
	RemoveFilm(ctx context.Context, filmID int, userID uuid.UUID) error
	GetAllFavourites(ctx context.Context, userID uuid.UUID) ([]models.DTOAllFavorites, error)
	//GetFilmByID(ctx context.Context, filmID int) (models.ResponseFilm, error)
	DeleteComment(ctx context.Context, ID uuid.UUID, userID uuid.UUID) error
	GetTwentyFilm(ctx context.Context) ([]models.DTOFilmMain, error)
	GetAllGenresWithCount(ctx context.Context) (map[string]int, error)
	GetFilmByGenre(ctx context.Context, genre string, page int) ([]models.ResponseFilmGenre, error)
	GetFilmByName(ctx context.Context, name string) ([]models.ResponseFilmName, error)
	GetCommentsByFilmID(ctx context.Context, filmID int) ([]models.ResponseComments, error)
	GetFilmById(ctx context.Context, filmID int, filmsID []int) (models.ResponseFilmPage, error)
}
type DB interface {
	User() UserStorage
	Film() FilmStorage
}
