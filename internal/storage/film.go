package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"go.uber.org/zap"
)

var _FilmStorage = (*film)(nil)

type film struct {
	logger *zap.Logger
	pgx    *pgxpool.Pool
}

func NewFilm(logger *zap.Logger, pgxPool *pgxpool.Pool) (*film, error) {
	film := &film{
		logger: logger,
		pgx:    pgxPool,
	}
	return film, nil
}

func (f *film) AddNewComment(ctx context.Context, data *models.DTOComments, userID uuid.UUID) error {
	const getData = `SELECT vote_average, vote_count FROM movie WHERE id = $1`
	const updateData = `UPDATE movie SET vote_average = $1, vote_count = $2 WHERE id = $3`
	const addComment = `INSERT INTO comments(id, user_id, film_id, comment) VALUES($1, $2, $3, $4)`
	tx, err := f.pgx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	var rating float64
	var count int
	err = tx.QueryRow(ctx, getData, data.FilmID).Scan(&rating, &count)
	if err != nil {
		return err
	}
	summa := average(rating, count)
	_, err = tx.Exec(ctx, updateData, summa, count+1, data.FilmID)
	if err != nil {
		return err
	}
	if len(data.Comment) > 0 {
		_, err = tx.Exec(ctx, addComment, data.ID, userID, data.FilmID, data.Comment)
		if err != nil {
			return err
		}
	}

	return nil
}

func average(rating float64, count int) float64 {
	return (rating*float64(count) + 1) / (float64(count) + 1)
}

func (f *film) AddFilmToFavourites(ctx context.Context, data *models.DTOFavorites) error {
	const addFilm = `INSERT INTO favorites(id, film_id, user_id) VALUES($1, $2, $3)`
	_, err := f.pgx.Exec(ctx, addFilm, data.ID, data.FilmID, data.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (f *film) RemoveFilm(ctx context.Context, filmID int, userID uuid.UUID) error {
	const removeFilm = `DELETE FROM favorites WHERE id = $1 AND user_id = $2`
	_, err := f.pgx.Exec(ctx, removeFilm, filmID, userID)
	if err != nil {
		return err
	}
	return nil
}
func (f *film) GetAllFavourites(ctx context.Context, userID uuid.UUID) ([]models.DTOAllFavorites, error) {
	const getAllFavourites = `
		SELECT f.film_id, m.title, m.weight_rating
		FROM favorites f
		JOIN movie m ON f.film_id = m.id
		WHERE f.user_id = $1
	`
	var favourites []models.DTOAllFavorites
	rows, err := f.pgx.Query(ctx, getAllFavourites, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp models.DTOAllFavorites
		err = rows.Scan(&temp.FilmID, &temp.Name, &temp.Rating)
		if err != nil {
			return nil, err
		}
		favourites = append(favourites, temp)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return favourites, nil
}

func (f *film) GetFilmByID(ctx context.Context, filmID int) (models.ResponseFilm, error) {
	const getFilm = `
		SELECT 
			m.title,
			m.genres,
			m.overview,
			m.production_companies,
			m.production_countries,
			m.release_date,
			m.run_time,
			m.vote_average,
			m.vote_count,
			m.director,
			m.weight_rating,
			a.name AS actor,
			k.keyword AS keyword,
			c.com AS comment,
			u.login AS user_login
		FROM movie m
		LEFT JOIN comments c ON m.id = c.film_id
		LEFT JOIN users u ON c.user_id = u.id
		LEFT JOIN movie_actor ma ON m.id = ma.movie_id
		LEFT JOIN actor a ON ma.actor_id = a.id
		LEFT JOIN movie_keywords mk ON m.id = mk.movie_id
		LEFT JOIN keywords k ON mk.keyword_id = k.id
		WHERE m.id = $1
		ORDER BY m.id, actor, keyword, comment, user_login;
	`

	var movie models.ResponseFilm
	var comments []string
	var userLogins []string
	var genres []string
	var actors []string
	var keyWords []string

	rows, err := f.pgx.Query(ctx, getFilm, filmID)
	if err != nil {
		return models.ResponseFilm{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var actor, keyword, comment, userLogin, genre string

		err = rows.Scan(
			&movie.Title,
			&genres,
			&movie.Overview,
			&movie.ProductionCompanies,
			&movie.ProductionCountries,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.VoteAverage,
			&movie.VoteCount,
			&movie.Director,
			&movie.WeightRating,
			&actor,
			&keyword,
			&comment,
			&userLogin,
		)
		if err != nil {
			return models.ResponseFilm{}, err
		}
		genres = append(genres, genre)
		actors = append(actors, actor)
		keyWords = append(keyWords, keyword)
		comments = append(comments, comment)
		userLogins = append(userLogins, userLogin)
	}

	for i := range comments {
		movie.FilmsComments = append(movie.FilmsComments, models.Comments{
			Com:       comments[i],
			UserLogin: userLogins[i],
		})
	}

	movie.Actor = actors
	movie.KeyWords = keyWords
	movie.Genres = genres

	return movie, nil
}
