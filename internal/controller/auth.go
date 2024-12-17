package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/recommender-system-for-MTUCI/2.0backend/internal/models"
	"go.uber.org/zap"
	"net/http"
)

// handle to registration, return jwt
func (ctrl *Controller) handleRegistration(ctx echo.Context) error {
	var req models.RequestRegister
	err := ctx.Bind(&req)
	if err != nil {
		ctrl.logger.Error("failed to bind registration request")
		return ctx.JSON(http.StatusBadRequest, err)
	}
	ctrl.logger.Info("received registration request")
	err = validateRegistration(req)
	if err != nil {
		ctrl.logger.Error("failed to validate registration request")
		return ctx.JSON(http.StatusBadRequest, err)
	}
	ctrl.logger.Info("validating registration request")
	password, err := hashPassword(req.Password)
	if err != nil {
		ctrl.logger.Error("failed to hash password")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully hashed password")
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
	ctrl.logger.Info("successfully sent registration response")
	err = ctrl.storage.User().AddUserInDB(ctrl.ctx, user)
	if err != nil {
		ctrl.logger.Error("failed to add user to database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully added user to database")
	accessToken, refreshToken, err := ctrl.generateAccessAndRefreshToken(user.ID)
	if err != nil {
		ctrl.logger.Error("failed to generate access token")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully generated access token")
	res := models.ResponseRegister{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return ctx.JSON(http.StatusOK, res)

}

// handle to log in, return jwt
func (ctrl *Controller) handleLogin(ctx echo.Context) error {
	var req models.RequestLogin
	err := ctx.Bind(&req)
	if err != nil {
		log.Error("failed to bind login request", err)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	ctrl.logger.Info("received login request")
	err = validateLogin(req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	id, password, err := ctrl.storage.User().GetUserIdByEmail(ctrl.ctx, req.Login)
	if err != nil {
		ctrl.logger.Error("failed to find user from database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully got user from database")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, id)
	if err != nil {
		ctrl.logger.Error("failed to find user from database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	ctrl.logger.Info("successfully got user from database")
	err = encryptPassword(req.Password, password)
	ctrl.logger.Info("successfully encrypted password")
	accessToken, refreshToken, err := ctrl.generateAccessAndRefreshToken(id)
	if err != nil {
		ctrl.logger.Error("failed to generate access token")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully generated access token")
	res := &models.ResponseLogin{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return ctx.JSON(http.StatusOK, res)
}

// handle to accept email
func (ctrl *Controller) handleAcceptEmail(ctx echo.Context) error {
	var req *models.RequestAcceptEmail
	id, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("failed to find user from database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully got user id")
	err = ctx.Bind(&req)
	if err != nil {
		ctrl.logger.Error("failed to bind accept email request")
		return ctx.JSON(http.StatusBadRequest, zap.Error(err))
	}
	ctrl.logger.Info("received accept email request")
	code, err := ctrl.storage.User().GetCodeFromDB(ctrl.ctx, id)
	if err != nil {
		ctrl.logger.Error("failed to find user from database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully got user from database")
	codeFromClient := req.Code
	if code != codeFromClient {
		ctrl.logger.Error("failed to accept email code from client")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully accept email code from client")
	err = ctrl.storage.User().UpdateUserStatus(ctrl.ctx, id)
	if err != nil {
		ctrl.logger.Error("failed to update user status")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully updated user status")
	return ctx.NoContent(http.StatusOK)
}

// handle to delete user
func (ctrl *Controller) handleDeleteUser(ctx echo.Context) error {
	id, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("failed to find user")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully got user id")
	err = ctrl.storage.User().DeleteUser(ctrl.ctx, id)
	if err != nil {
		ctrl.logger.Error("failed to delete user")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully deleted user")
	return ctx.NoContent(http.StatusOK)
}

// handle to change password
func (ctrl *Controller) handleChangePassword(ctx echo.Context) error {
	var req models.RequestChangePassword
	id, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("failed to parse uuid")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully got user id")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, id)
	if err != nil {
		ctrl.logger.Error("failed to find user from database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "email dont accept"})
	}
	ctrl.logger.Info("successfully got user status from database")
	err = ctx.Bind(&req)
	if err != nil {
		ctrl.logger.Error("failed to bind change password request")
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	ctrl.logger.Info("received change password request")
	err = validatePassword(req)
	if err != nil {
		ctrl.logger.Error("failed to validate change password request")
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	ctrl.logger.Info("successfully validate change password request")
	oldPassword, err := ctrl.storage.User().GetPasswordFromDB(ctrl.ctx, id)
	if err != nil {
		ctrl.logger.Error("failed to retrieve password from user")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully got old password from user")
	err = encryptPassword(req.OldPassword, oldPassword)
	if err != nil {
		ctrl.logger.Error("failed to encrypt old password")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully encrypted old password")
	hashPas, err := hashPassword(req.NewPassword)
	if err != nil {
		ctrl.logger.Error("failed to hash password")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully hashed password")
	data := &models.DTOPassword{
		hashPas,
	}
	err = ctrl.storage.User().UpdatePassword(ctrl.ctx, id, data.Password)
	if err != nil {
		ctrl.logger.Error("failed to update password")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully updated password")
	return ctx.NoContent(http.StatusNoContent)
}

// handle to get user page
func (ctrl *Controller) handleGetMe(ctx echo.Context) error {
	userID, err := ctrl.getUserId(ctx.Request())
	if err != nil {
		ctrl.logger.Error("failed to find user from database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully got user id")
	status, err := ctrl.storage.User().GetStatusFromUser(ctrl.ctx, userID)
	if err != nil {
		ctrl.logger.Error("failed to find user from database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if status == false {
		ctrl.logger.Error("email dont accept")
		return ctx.JSON(http.StatusForbidden, err)
	}
	ctrl.logger.Info("successfully got user status from database")
	login, err := ctrl.storage.User().GetMe(ctrl.ctx, userID)
	if err != nil {
		ctrl.logger.Error("failed to find user from database")
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	ctrl.logger.Info("successfully got user status from database")
	return ctx.JSON(http.StatusOK, echo.Map{"login": login})
}

// handle to update tokens
func (ctrl *Controller) handleRefresh(ctx echo.Context) error {
	refreshToken, err := tokenInHeader(ctx.Request())
	if err != nil {
		ctrl.logger.Error("failed to get refresh token from header", zap.Error(err))
		return ctx.JSON(http.StatusBadRequest, err)
	}
	ctrl.logger.Info("successfully got refresh token from header")

	userID, err := ctrl.validateRefreshToken(refreshToken)
	if err != nil {
		ctrl.logger.Error("failed to validate refresh token", zap.Error(err))
		return ctx.JSON(http.StatusUnauthorized, zap.Error(err))
	}
	ctrl.logger.Info("successfully got refresh token from header")

	newAccessToken, newRefreshToken, err := ctrl.generateAccessAndRefreshToken(userID)
	if err != nil {
		ctrl.logger.Error("failed to generate new tokens", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate tokens"})
	}

	ctrl.logger.Info("successfully generated new tokens")

	response := models.ResponseRegister{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}
	ctrl.logger.Info("successfully refreshed tokens")
	return ctx.JSON(http.StatusOK, response)
}
