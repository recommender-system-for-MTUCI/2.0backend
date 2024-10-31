package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"net/http"
)

func (ctrl *Controller) handleRegistration(ctx echo.Context) error {
	var req models.RequestRegister
	err := ctx.Bind(&req)
	if err != nil {
		ctrl.logger.Error("failed to bind registration request")
		return ctx.JSON(http.StatusBadRequest, err)
	}
	err = validateRegistration(req)
	if err != nil {
		ctrl.logger.Error("failed to validate registration request")
		return ctx.JSON(http.StatusBadRequest, err)
	}
	password, err := hashPassword(req.Password)
	if err != nil {
		ctrl.logger.Error("failed to hash password")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	user := models.DTORegister{
		ID:       uuid.New(),
		Login:    req.Login,
		Password: password,
	}
	err = ctrl.sendMessages(user.Login)
	if err != nil {
		ctrl.logger.Error("failed to send registration response")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	///need add logic for insert user in db
	accessToken, refreshToken, err := ctrl.generateAccessAndRefreshToken(user.ID)
	if err != nil {
		ctrl.logger.Error("failed to generate access token")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	res := models.ResponseRegister{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return ctx.JSON(http.StatusOK, res)

}
func (ctrl *Controller) handleLogin(ctx echo.Context, req *models.RequestLogin) error {
	err := ctx.Bind(&req)
	if err != nil {
		log.Error("failed to bind login request", err)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	err = validateLogin(req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{})
}

func (ctrl *Controller) handleAcceptEmail(ctx echo.Context) error {
	return nil
}
func (ctrl *Controller) handleLogout(ctx echo.Context, req *models.RequestLogin) error {
	return nil
}

func (ctrl *Controller) handleDelete(ctx echo.Context, req *models.RequestLogin) error {
	return nil
}
