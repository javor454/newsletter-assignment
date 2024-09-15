package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/javor454/newsletter-assignment/app/http-server"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/ui/http/request"
)

type RegisterUserHandler interface {
	Handle(ctx context.Context, email string, password string) (*dto.Token, error)
}

type LoginUserHandler interface {
	Handle(ctx context.Context, email string, password string) (*dto.Token, error)
}

type UserController struct {
	lg  logger.Logger
	ruh RegisterUserHandler
	luh LoginUserHandler
}

func NewUserController(
	lg logger.Logger,
	httpServer *http_server.HttpServer,
	ruh RegisterUserHandler,
	luh LoginUserHandler,
) *UserController {
	u := &UserController{
		ruh: ruh,
		luh: luh,
		lg:  lg,
	}

	httpServer.GetEngine().POST("api/users/register", u.RegisterUser)
	httpServer.GetEngine().POST("api/users/login", u.LoginUser)

	return u
}

func (u *UserController) RegisterUser(ctx *gin.Context) {
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

	var req *request.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.lg.WithError(err).Error("Failed to bind request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	token, err := u.ruh.Handle(ctx, req.Email, req.Password)
	if err != nil {
		u.lg.WithError(err).Error("Failed to handle register")
		ctx.JSON(http.StatusInternalServerError, gin.H{})

		return
	}

	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", token.String()))
	ctx.JSON(http.StatusCreated, gin.H{})
}

func (u *UserController) LoginUser(ctx *gin.Context) {
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

	var req *request.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.lg.WithError(err).Error("Failed to bind request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	token, err := u.luh.Handle(ctx, req.Email, req.Password)
	if err != nil {
		u.lg.WithError(err).Error("Failed to handle login")
		ctx.JSON(http.StatusInternalServerError, gin.H{})

		return
	}

	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", token.String()))
	ctx.JSON(http.StatusCreated, gin.H{})
}
