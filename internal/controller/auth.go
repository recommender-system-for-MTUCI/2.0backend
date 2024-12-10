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
		Code:         generationRandomCode(),
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
	id, password, err := ctrl.storage.User().GetUserIdByEmail(ctrl.ctx, req.Login)
	if err != nil {
		ctrl.logger.Error("failed to find user from database")
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
	err = validateLogin(req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	err = encryptPassword(req.Password, password)
	accessToken, refreshToken, err := ctrl.generateAccessAndRefreshToken(id)
	if err != nil {
		ctrl.logger.Error("failed to generate access token")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	res := &models.ResponseLogin{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return ctx.JSON(http.StatusOK, res)
}

func (ctrl *Controller) handleAcceptEmail(ctx echo.Context) error {
	id, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("failed to find user from database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	var req *models.RequestAcceptEmail
	err = ctx.Bind(&req)
	if err != nil {
		log.Error("failed to bind accept email request", err)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
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
	err = ctrl.storage.User().UpdateUserStatus(ctrl.ctx, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, echo.Map{})
}

func (ctrl *Controller) handleDelete(ctx echo.Context) error {
	id, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("failed to find user")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	err = ctrl.storage.User().DeleteUser(ctrl.ctx, id)
	if err != nil {
		ctrl.logger.Error("failed to delete user")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusNoContent, echo.Map{})
}
func (ctrl *Controller) handleChangePassword(ctx echo.Context) error {
	id, err := ctrl.getUserId(ctx.Request())
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
	err = encryptPassword(req.OldPassword, oldPassword)
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
	return ctx.NoContent(http.StatusNoContent)
}

func (ctrl *Controller) handleGetMe(ctx echo.Context) error {
	userID, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("failed to find user from database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	login, err := ctrl.storage.User().GetMe(ctrl.ctx, userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, echo.Map{"login": login})

}
