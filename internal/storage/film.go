package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"go.uber.org/zap"
)

var _FilmStorage = (*film)(nil)

type film struct {
	logger *zap.Logger
	pgx    *pgxpool.Pool
}

// create new instanse film storage
func NewFilm(logger *zap.Logger, pgxPool *pgxpool.Pool) (*film, error) {
	film := &film{
		logger: logger,
		pgx:    pgxPool,
	}
	return film, nil
}

func (f *film) AddNewComment(ctx context.Context, data *models.DTOComments) error {
	const getData = `SELECT vote_average, vote_count FROM movie WHERE id = $1`
	const updateData = `UPDATE movie SET vote_average = $1, vote_count = $2 WHERE id = $3`
	const addComment = `INSERT INTO comments(id, film_id, user_id, comment, rating, created_at) VALUES($1, $2, $3, $4, $5, NOW())`

	tx, err := f.pgx.Begin(ctx)
	if err != nil {
		f.logger.Error("failed to start transaction", zap.Error(err))
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			f.logger.Error("panic", zap.Any("panic", p))
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			f.logger.Error("rollback transaction", zap.Any("err", err))
			_ = tx.Rollback(ctx)
		} else {
			f.logger.Info("commit transaction")
			err = tx.Commit(ctx)
		}
	}()

	var rating float64
	var count int
	err = tx.QueryRow(ctx, getData, data.FilmID).Scan(&rating, &count)
	if err != nil {
		f.logger.Error("failed to get film data", zap.Any("err", err))
		return err
	}
	summa := average(rating, count)
	_, err = tx.Exec(ctx, updateData, summa, count+1, data.FilmID)
	if err != nil {
		f.logger.Error("failed to update film data", zap.Any("err", err))
		return err
	}
	_, err = tx.Exec(ctx, addComment, data.ID, data.FilmID, data.UserID, data.Comment, data.Rating)
	if err != nil {
		f.logger.Error("failed to add comment", zap.Any("err", err))
		return err
	}
	f.logger.Info("successfully added comment", zap.Any("comment", data.Comment))
	return nil
}

// func to update films rating
func average(rating float64, count int) float64 {
	return (rating*float64(count) + 1) / (float64(count) + 1)
}

// add film to favourites
func (f *film) AddFilmToFavourites(ctx context.Context, data *models.DTOFavorites) error {
	const checkFilm = `SELECT id FROM favorites WHERE film_id = $1 AND user_id = $2`
	const addFilm = `INSERT INTO favorites(id, film_id, user_id) VALUES($1, $2, $3)`
	var existingID uuid.UUID
	err := f.pgx.QueryRow(ctx, checkFilm, data.FilmID, data.UserID).Scan(&existingID)
	if err == nil {
		f.logger.Error("failed to add in favourites, film already in db", zap.Error(err))
		return err
	} else if err != pgx.ErrNoRows {
		f.logger.Error("failed to get in db", zap.Error(err))
		return err
	}
	_, err = f.pgx.Exec(ctx, addFilm, data.ID, data.FilmID, data.UserID)
	if err != nil {
		f.logger.Error("failed to add in favourites", zap.Error(err))
		return err
	}
	f.logger.Info("film successfuuly added in db")
	return nil
}

// func to remove film from favouritees
func (f *film) RemoveFilm(ctx context.Context, filmID int, userID uuid.UUID) error {
	const removeFilm = `DELETE FROM favorites WHERE film_id = $1 AND user_id = $2`
	_, err := f.pgx.Exec(ctx, removeFilm, filmID, userID)
	if err != nil {
		f.logger.Error("failed to remove favorites", zap.Error(err))
		return err
	}
	f.logger.Info("film successfuuly removed in db")
	return nil
}
func (f *film) GetAllFavourites(ctx context.Context, userID uuid.UUID) ([]models.DTOAllFavorites, error) {
	const getAllFavourites = `SELECT f.film_id, m.title, m.weight_rating
		FROM favorites f
		JOIN movie m ON f.film_id = m.id
		WHERE f.user_id = $1
	`
	var favourites []models.DTOAllFavorites
	rows, err := f.pgx.Query(ctx, getAllFavourites, userID)
	if err != nil {
		f.logger.Error("failed to get all favourites", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var temp models.DTOAllFavorites
		err = rows.Scan(&temp.FilmID, &temp.Name, &temp.Rating)
		if err != nil {
			f.logger.Error("failed to scan favourites", zap.Error(err))
			return nil, err
		}
		favourites = append(favourites, temp)
	}
	f.logger.Info("successfully retrieved all favourites")
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return favourites, nil
}

// need change all logic
/*func (f *film) GetFilmByID(ctx context.Context, filmID int) (models.ResponseFilm, error) {
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
}*/

// func to delete comment
func (f *film) DeleteComment(ctx context.Context, ID uuid.UUID, userID uuid.UUID) error {
	const delete = `DELETE FROM comments WHERE id = $1 AND user_id = $2;`
	_, err := f.pgx.Exec(ctx, delete, ID, userID)
	if err != nil {
		f.logger.Error("failed to delete comment", zap.Error(err))
		return err
	}
	f.logger.Info("successfully deleted comment")
	return nil
}

// func to get 20 most popular film to main page
func (f *film) GetTwentyFilm(ctx context.Context) ([]models.DTOFilmMain, error) {
	const get = `SELECT id, title, vote_average FROM movie ORDER BY vote_average DESC, release_data DESC LIMIT 20;`
	var films []models.DTOFilmMain
	rows, err := f.pgx.Query(ctx, get)
	if err != nil {
		f.logger.Error("failed to get all films", zap.Error(err))
		return films, err
	}
	defer rows.Close()
	for rows.Next() {
		var movie models.DTOFilmMain
		err = rows.Scan(&movie.FilmID, &movie.Name, &movie.Rating)
		if err != nil {
			f.logger.Error("failed to scan", zap.Error(err))
			return films, err
		}
		films = append(films, movie)
	}
	f.logger.Info("successfully retrieved all films")
	return films, nil
}

// func to get all genres with count film for each genres
func (f *film) GetAllGenresWithCount(ctx context.Context) (map[string]int, error) {
	const getGenresWithCount = `
		SELECT unnest(genres) AS genre, COUNT(*) AS film_count
		FROM movie
		GROUP BY genre;
	`

	genres := make(map[string]int)

	rows, err := f.pgx.Query(ctx, getGenresWithCount)
	if err != nil {
		f.logger.Error("failed to query genres with count", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var genre string
		var count int
		err = rows.Scan(&genre, &count)
		if err != nil {
			f.logger.Error("failed to scan row", zap.Error(err))
			return nil, err
		}
		genres[genre] = count
	}
	f.logger.Info("successfully retrieved all genres with count")
	return genres, nil
}

// func to get films by genre
func (f *film) GetFilmByGenre(ctx context.Context, genre string, page int) ([]models.ResponseFilmGenre, error) {
	const pageSize = 30
	const getFilm = `
		SELECT id, title, vote_average 
		FROM movie 
		WHERE $1 = ANY (genres) 
		LIMIT $2 OFFSET $3;
	`

	var films []models.ResponseFilmGenre
	offset := (page - 1) * pageSize
	rows, err := f.pgx.Query(ctx, getFilm, genre, pageSize, offset)
	if err != nil {
		f.logger.Error("failed to query film", zap.Error(err))
		return films, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.ResponseFilmGenre
		err = rows.Scan(&movie.FilmID, &movie.Name, &movie.Rating)
		if err != nil {
			f.logger.Error("failed to scan", zap.Error(err))
			return films, err
		}
		films = append(films, movie)
	}
	f.logger.Info("successfully retrieved all films")

	return films, nil
}

// func to get film by search
func (f *film) GetFilmByName(ctx context.Context, name string) ([]models.ResponseFilmName, error) {
	const getFilm = `SELECT id, title, vote_average FROM movie WHERE title ILIKE $1`
	searchName := "%" + name + "%"
	var films []models.ResponseFilmName

	rows, err := f.pgx.Query(ctx, getFilm, searchName)
	if err != nil {
		f.logger.Error("failed to get film", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var movie models.ResponseFilmName
		if err := rows.Scan(&movie.FilmID, &movie.Name, &movie.Rating); err != nil {
			f.logger.Error("failed to scan film", zap.Error(err))
			return nil, err
		}
		films = append(films, movie)
	}
	f.logger.Info("successfully retrieved all films")
	return films, nil
}

// func to get all comments by film id
func (f *film) GetCommentsByFilmID(ctx context.Context, filmID int) ([]models.ResponseComments, error) {
	const getAllComments = `
		SELECT c.id, c.rating,u.login, c.comment
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.film_id = $1
		ORDER BY c.created_at DESC`

	rows, err := f.pgx.Query(ctx, getAllComments, filmID)
	if err != nil {
		f.logger.Error("failed to get comments", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var comments []models.ResponseComments
	for rows.Next() {
		var comment models.ResponseComments
		if err := rows.Scan(&comment.ID, &comment.Rating, &comment.UserLogin, &comment.Comment); err != nil {
			f.logger.Error("failed to scan comment", zap.Error(err))
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		f.logger.Error("rows iteration error", zap.Error(err))
		return nil, err
	}
	f.logger.Info("successfully retrieved all comments")
	return comments, nil
}
