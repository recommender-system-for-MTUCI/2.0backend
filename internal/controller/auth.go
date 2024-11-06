package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"net/http"
)

func (ctrl *Controller) handleRegistration(ctx echo.Context) error {
	var req *models.RequestRegister
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
	user := &models.DTORegister{
		ID:           uuid.New(),
		Login:        req.Login,
		Password:     password,
		Confirmation: false,
		Code:         ctrl.generationRandomCode(),
	}
	err = ctrl.sendMessages(user.Login, user.Code)
	if err != nil {
		ctrl.logger.Error("failed to send registration response")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	err = ctrl.storage.User().AddUserInDB(ctrl.ctx, user)
	if err != nil {
		ctrl.logger.Error("failed to add user to database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
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
func (ctrl *Controller) handleLogin(ctx echo.Context) error {
	var req *models.RequestLogin
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
	//need add logic for request id
	var req *models.RequestAcceptEmail
	err := ctx.Bind(&req)
	if err != nil {
		log.Error("failed to bind accept email request", err)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	id, err := uuid.Parse("766c9d59-a163-474d-a146-dc2d69cdfa40")
	if err != nil {
		ctrl.logger.Error("rjfekfkr")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	code, err := ctrl.storage.User().GetCodeFromDB(ctrl.ctx, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	codeFromClient := req.Code
	if code != codeFromClient {
		ctrl.logger.Error("failed to accept email code from client")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, echo.Map{})
}
func (ctrl *Controller) handleLogout(ctx echo.Context, req *models.RequestLogin) error {
	return nil
}

func (ctrl *Controller) handleDelete(ctx echo.Context) error {
	return nil
}
func (ctrl *Controller) handleChangePassword(ctx echo.Context) error {
	id, err := uuid.Parse("766c9d59-a163-474d-a146-dc2d69cdfa40")
	if err != nil {
		ctrl.logger.Error("failed to parse uuid")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	var req *models.RequestChangePassword
	err = ctx.Bind(&req)
	if err != nil {
		ctrl.logger.Error("failed to bind change password request")
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	err = validatePassword(req)
	if err != nil {
		ctrl.logger.Error("failed to validate change password request")
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	oldPassword, err := ctrl.storage.User().GetPasswordFromDB(ctrl.ctx, id)
	if err != nil {
		ctrl.logger.Error("failed to retrieve password from user")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if req.OldPassword != oldPassword {
		ctrl.logger.Error("invalid old password")
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid old password"})
	}
	hashPas, err := hashPassword(req.NewPassword)
	if err != nil {
		ctrl.logger.Error("failed to hash password")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	pas := &models.DTOPassword{
		hashPas,
	}
	err = ctrl.storage.User().UpdatePassword(ctrl.ctx, id, pas.Password)
	if err != nil {
		ctrl.logger.Error("failed to update password")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, echo.Map{})
}
