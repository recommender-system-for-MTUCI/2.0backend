package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"net/http"
)

func (ctrl *Controller) handleRegistration(ctx echo.Context, req *models.RequestLogin) error {
	_ = ctx.Bind(&req)
	return nil
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
func (ctrl *Controller) handleLogout(ctx echo.Context, req *models.RequestLogin) error {
	return nil
}

func (ctrl *Controller) handleDelete(ctx echo.Context, req *models.RequestLogin) error {
	return nil
}
