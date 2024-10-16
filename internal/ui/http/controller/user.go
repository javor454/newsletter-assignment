package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/javor454/newsletter-assignment/app/http_server"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/internal/application"
	"github.com/javor454/newsletter-assignment/internal/ui/http/request"
)

type RegisterUserHandler interface {
	Handle(ctx context.Context, email string, password string) (string, error)
}

type LoginUserHandler interface {
	Handle(ctx context.Context, email string, password string) (string, error)
}

type UserController struct {
	lg  logger.Logger
	ruh RegisterUserHandler
	luh LoginUserHandler
}

func NewUserController(
	lg logger.Logger,
	ruh RegisterUserHandler,
	luh LoginUserHandler,
) *UserController {
	controller := &UserController{
		ruh: ruh,
		luh: luh,
		lg:  lg,
	}

	return controller
}

func (u *UserController) RegisterUserController(httpServer *http_server.Server) {
	httpServer.GetEngine().POST("api/v1/users/register", u.Register)
	httpServer.GetEngine().POST("api/v1/users/login", u.Login)
}

// Register
//
//	@Summary	Register user, returning token for authorization
//	@Router		/api/v1/users/register [post]
//	@Tags		public user
//	@Accepts	json
//	@Produce	json
//
//	@Param		data	body	request.UserRequest	true	"Data for registering user"
//
//	@Success	201		"User was successfully registered"
//	@Failure	400		{object}	response.Error	"Invalid request with detail"
//	@Failure	409		{object}	response.Error	"Email taken"
//	@Failure	500		"Unexpected exception"
func (u *UserController) Register(ctx *gin.Context) {
	var h *request.ContentTypeHeader
	if err := ctx.ShouldBindHeader(&h); err != nil {
		u.lg.WithError(err).Error("Failed to bind headers")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}
	if err := h.Validate(); err != nil {
		u.lg.WithError(err).Error("Failed to validate headers")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	var req *request.UserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.lg.WithError(err).Error("Failed to bind request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	token, err := u.ruh.Handle(ctx, req.Email, req.Password)
	if err != nil {
		code, body := func(err error) (int, gin.H) {
			if errors.Is(err, application.EmailTakenError) {
				return http.StatusConflict, gin.H{"error": "Email taken"}
			}

			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to handle register")
		ctx.JSON(code, body)

		return
	}

	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	ctx.JSON(http.StatusCreated, gin.H{})
}

// Login
//
//	@Summary	Login user, returning token for authorization
//	@Router		/api/v1/users/login [post]
//	@Tags		public user
//	@Accepts	json
//	@Produce	json
//
//	@Param		data	body	request.UserRequest	true	"Data for user login"
//
//	@Success	201		"User successfully logged in"
//	@Failure	400		{object}	response.Error	"Invalid request with detail"
//	@Failure	401		{object}	response.Error	"Invalid credentials"
//	@Failure	500		"Unexpected exception"
func (u *UserController) Login(ctx *gin.Context) {
	var h *request.ContentTypeHeader
	if err := ctx.ShouldBindHeader(&h); err != nil {
		u.lg.WithError(err).Error("Failed to bind headers")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}
	if err := h.Validate(); err != nil {
		u.lg.WithError(err).Error("Failed to validate headers")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	var req *request.UserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.lg.WithError(err).Error("Failed to bind request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	token, err := u.luh.Handle(ctx, req.Email, req.Password)
	if err != nil {
		code, body := func(err error) (int, gin.H) {
			if errors.Is(err, application.UserNotFoundError) || errors.Is(err, application.InvalidPasswordError) {
				return http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}
			}

			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to handle login")
		ctx.JSON(code, body)

		return
	}

	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	ctx.JSON(http.StatusCreated, gin.H{})
}
