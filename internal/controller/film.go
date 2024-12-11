package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (ctrl *Controller) handleAddComments(ctx echo.Context) error {
	var req *models.Comment
	err := ctx.Bind(&req)
	if err != nil {
		ctrl.logger.Error("got err while bind request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	id, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, id)
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
	DTO := &models.DTOComments{
		ID:      uuid.New(),
		FilmID:  filmID,
		Comment: req.Comment,
	}
	err = ctrl.storage.Film().AddNewComment(ctrl.ctx, DTO, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.NoContent(http.StatusCreated)
}

func (ctrl *Controller) handleDeleteComments(ctx echo.Context) error {
	return nil
}
func (ctrl *Controller) handleAddToFavorites(ctx echo.Context) error {
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
	strID := ctx.Param("id")
	//ctrl.logger.Error("strID", zap.Any("strID", strID))
	if strID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing filmID")
	}
	filmID, err := strconv.Atoi(strID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid film ID format")
	}
	data := &models.DTOFavorites{
		ID:     uuid.New(),
		FilmID: filmID,
		UserID: userID,
	}
	err = ctrl.storage.Film().AddFilmToFavourites(ctrl.ctx, data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.NoContent(http.StatusCreated)
}

func (ctrl *Controller) handleDeleteFromFavorites(ctx echo.Context) error {
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
	err = ctrl.storage.Film().RemoveFilm(ctrl.ctx, filmID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.NoContent(http.StatusCreated)
}

// need to add rating
func (ctrl *Controller) handleGetFavorites(ctx echo.Context) error {
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
	data, err := ctrl.storage.Film().GetAllFavourites(ctrl.ctx, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if len(data) == 0 {
		return ctx.NoContent(http.StatusNoContent)
	}

	return ctx.JSON(http.StatusOK, data)
}

// need add
func (ctrl *Controller) handleGetMainPage(ctx echo.Context) error {
	//here will be work with client ml
	return ctx.HTML(http.StatusOK, "main.html")
}
func (ctrl *Controller) handleGetFilmByID(ctx echo.Context) error {
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
}
