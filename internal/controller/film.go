package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// handle to add comment to film
func (ctrl *Controller) handleAddComment(ctx echo.Context) error {
	var req *models.Comment
	err := ctx.Bind(&req)
	if err != nil {
		ctrl.logger.Error("got err while bind request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully bind request")
	id, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("got err while get user id", zap.Error(err))
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	ctrl.logger.Info("successfully get userID")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, id)
	if err != nil {
		ctrl.logger.Error("got err while get user status", zap.Error(err))
		return ctx.JSON(http.StatusForbidden, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	ctrl.logger.Info("successfully get user status")
	filmID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctrl.logger.Error("got err while convert film id to int", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully convert film id")
	DTO := &models.DTOComments{
		ID:      uuid.New(),
		FilmID:  filmID,
		UserID:  id,
		Comment: req.Comment,
		Rating:  req.Rating,
	}
	err = ctrl.storage.Film().AddNewComment(ctrl.ctx, DTO)
	if err != nil {
		ctrl.logger.Error("got err while add new comment", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully add new comment")
	return ctx.NoContent(http.StatusCreated)
}

// handle to delete comment
func (ctrl *Controller) handleDeleteComment(ctx echo.Context) error {
	userID, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("got err while get user id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully get user id")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, userID)
	if err != nil {
		ctrl.logger.Error("got err while get user status", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	ctrl.logger.Info("successfully get user status")
	commentID := ctx.Param("id")
	if commentID == "" {
		ctrl.logger.Error("got empty comment id")
		return echo.NewHTTPError(http.StatusBadRequest, "missing filmID")
	}
	ID, err := uuid.Parse(commentID)
	if err != nil {
		ctrl.logger.Error("got err while convert comment id to int", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully convert comment id")
	err = ctrl.storage.Film().DeleteComment(ctrl.ctx, ID, userID)
	if err != nil {
		ctrl.logger.Error("got err while delete comment", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully delete comment")
	return ctx.NoContent(http.StatusNoContent)
}

// handle for add film in favourites
func (ctrl *Controller) handleAddToFavourites(ctx echo.Context) error {
	userID, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("got err while get user id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully get user id")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, userID)
	if err != nil {
		ctrl.logger.Error("got err while get user status", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	ctrl.logger.Info("successfully get user status")
	strID := ctx.Param("id")
	if strID == "" {
		ctrl.logger.Error("got empty filmID")
		return echo.NewHTTPError(http.StatusBadRequest, "missing filmID")
	}
	filmID, err := strconv.Atoi(strID)
	if err != nil {
		ctrl.logger.Error("got err while convert film id to int", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid film ID format")
	}
	data := &models.DTOFavorites{
		ID:     uuid.New(),
		FilmID: filmID,
		UserID: userID,
	}
	err = ctrl.storage.Film().AddFilmToFavourites(ctrl.ctx, data)
	if err != nil {
		ctrl.logger.Error("got err while add to favourites", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully add to favourites")
	return ctx.NoContent(http.StatusCreated)
}

func (ctrl *Controller) handleDeleteFromFavorites(ctx echo.Context) error {
	userID, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("got err while get user id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully get user id")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, userID)
	if err != nil {
		ctrl.logger.Error("got err while get user status", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	ctrl.logger.Info("successfully get user status")
	filmIDStr := ctx.Param("id")
	if filmIDStr == "" {
		ctrl.logger.Error("got empty filmID")
		return echo.NewHTTPError(http.StatusBadRequest, "missing filmID")
	}
	filmID, err := strconv.Atoi(filmIDStr)
	if err != nil {
		ctrl.logger.Error("got err while convert film id to int", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully convert film id")
	err = ctrl.storage.Film().RemoveFilm(ctrl.ctx, filmID, userID)
	if err != nil {
		ctrl.logger.Error("got err while delete favourites", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully delete favourites")
	return ctx.NoContent(http.StatusNoContent)
}

// handle to get all favourites by ID
func (ctrl *Controller) handleGetFavourites(ctx echo.Context) error {
	userID, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("got err while get user id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully get user id")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, userID)
	if err != nil {
		ctrl.logger.Error("got err while get user status", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	ctrl.logger.Info("successfully get user status")
	data, err := ctrl.storage.Film().GetAllFavourites(ctrl.ctx, userID)
	if err != nil {
		ctrl.logger.Error("got err while get favourites", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully get favourites")
	return ctx.JSON(http.StatusOK, data)
}

// handle to get main page
func (ctrl *Controller) handleGetMainPage(ctx echo.Context) error {
	data, err := ctrl.storage.Film().GetTwentyFilm(ctrl.ctx)
	if err != nil {
		ctrl.logger.Error("got err while get  film", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully get film")
	return ctx.JSON(http.StatusOK, data)
}

// need add logic with db and ml
/*func (ctrl *Controller) handleGetFilmByID(ctx echo.Context) error {
	userID, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	filmID, err := strconv.Atoi(ctx.Param("film_id"))
	if err != nil {
		ctrl.logger.Error("got err while convert film id to int", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	data, err := ctrl.storage.Film().GetFilmByID(ctrl.ctx, filmID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)

	}
	return ctx.JSON(http.StatusOK, data)
}*/

// handle for frontend developer, return all genres for filter
func (ctrl *Controller) handleGetAllGenres(ctx echo.Context) error {
	userID, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("got err while get user id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully get user id")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, userID)
	if err != nil {
		ctrl.logger.Error("got err while get user status", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	ctrl.logger.Info("successfully get user status")
	data, err := ctrl.storage.Film().GetAllGenresWithCount(ctrl.ctx)
	if err != nil {
		ctrl.logger.Error("got err while get genres", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully get genres")
	return ctx.JSON(http.StatusOK, data)
}

// handle to get films by genre
func (ctrl *Controller) handleGetFilmsByGenre(ctx echo.Context) error {
	userID, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("got err while get user id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully get user id")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, userID)
	if err != nil {
		ctrl.logger.Error("got err while get user status", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	ctrl.logger.Info("successfully get user status")
	genre := ctx.Param("genre")
	if genre == "" {
		ctrl.logger.Error("got empty genre")
		return echo.NewHTTPError(http.StatusBadRequest, "missing filmID")
	}
	ctrl.logger.Info("successfully get user status")
	pageStr := ctx.Param("page")
	if pageStr == "" {
		ctrl.logger.Error("got empty page")
		return echo.NewHTTPError(http.StatusBadRequest, "missing filmID")
	}
	ctrl.logger.Info("successfully get user status")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctrl.logger.Error("got err while convert page to int", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully get user status")
	data, err := ctrl.storage.Film().GetFilmByGenre(ctrl.ctx, genre, page)
	if err != nil {
		ctrl.logger.Error("got err while get film", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully get film")
	return ctx.JSON(http.StatusOK, data)
}

// handle to find film by name
func (ctrl *Controller) handleGetFilmsByName(ctx echo.Context) error {
	userID, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("got err while get user id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully get user id")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, userID)
	if err != nil {
		ctrl.logger.Error("got err while get user status", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	ctrl.logger.Info("successfully get user status")
	name := ctx.Param("name")
	if name == "" {
		ctrl.logger.Error("got empty name")
		return echo.NewHTTPError(http.StatusBadRequest, "missing filmID")
	}
	ctrl.logger.Info("successfully get user status")
	data, err := ctrl.storage.Film().GetFilmByName(ctrl.ctx, name)
	if err != nil {
		ctrl.logger.Error("got err while get film", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully get film")
	return ctx.JSON(http.StatusOK, data)
}

// handle to get all comments by filmID
func (ctrl *Controller) handleGetCommentsByFilmID(ctx echo.Context) error {
	userID, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("got err while get user id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully get user id")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, userID)
	if err != nil {
		ctrl.logger.Error("got err while get user status", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	ctrl.logger.Info("successfully get user status")
	IDStr := ctx.Param("id")
	if IDStr == "" {
		ctrl.logger.Error("got empty id")
		return echo.NewHTTPError(http.StatusBadRequest, "missing filmID")
	}
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		ctrl.logger.Error("got err while convert page to int", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctrl.logger.Info("successfully get film id")

	data, err := ctrl.storage.Film().GetCommentsByFilmID(ctrl.ctx, ID)
	if err != nil {
		ctrl.logger.Error("failed to get data from db", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	ctrl.logger.Info("successfully get comments")
	return ctx.JSON(http.StatusOK, data)

}
