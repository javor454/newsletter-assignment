package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/javor454/newsletter-assignment/app/http_server"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/internal/domain"
	"github.com/javor454/newsletter-assignment/internal/ui/http/middleware"
	"github.com/javor454/newsletter-assignment/internal/ui/http/request"
	"github.com/javor454/newsletter-assignment/internal/ui/http/response"
)

type CreateNewsletterHandler interface {
	Handle(ctx context.Context, userId, name string, description *string) error
}

type GetNewslettersByUserIDHandler interface {
	Handle(ctx context.Context, userID string, pageSize, pageNumber int) ([]*domain.Newsletter, error)
}

type NewsletterController struct {
	lg                     logger.Logger
	createNewsletter       CreateNewsletterHandler
	getNewslettersByUserID GetNewslettersByUserIDHandler
}

func NewNewsletterController(
	lg logger.Logger,
	httpServer *http_server.Server,
	cnh CreateNewsletterHandler,
	gnbui GetNewslettersByUserIDHandler,
	authMiddleware *middleware.AuthMiddleware,
) *NewsletterController {
	controller := &NewsletterController{createNewsletter: cnh, getNewslettersByUserID: gnbui, lg: lg}

	httpServer.GetEngine().POST("api/v1/newsletters/create", authMiddleware.Handle, controller.Create)
	httpServer.GetEngine().GET("api/v1/newsletters", authMiddleware.Handle, controller.GetNewslettersByUserID)

	return controller
}

func (u *NewsletterController) Create(ctx *gin.Context) {
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

	var req *request.CreateNewsletterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.lg.WithError(err).Error("Failed to bind request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	userID, ok := ctx.Get(middleware.UserIDKey)
	if !ok {
		u.lg.Error("User ID missing in gin context")
		ctx.JSON(http.StatusInternalServerError, gin.H{})

		return
	}

	if err := u.createNewsletter.Handle(ctx, userID.(string), req.Name, req.Description); err != nil {
		u.lg.WithError(err).Error("Failed to create newsletter")
		ctx.JSON(http.StatusInternalServerError, gin.H{})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func (u *NewsletterController) GetNewslettersByUserID(ctx *gin.Context) {
	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil {
		u.lg.WithError(err).Error("Failed to parse page size")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid pageSize"})

		return
	}
	if pageSize < 1 {
		u.lg.WithError(err).Error("Failed to parse page size")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid pageSize"})

		return
	}

	pageNumber, err := strconv.Atoi(ctx.DefaultQuery("page_number", "1"))
	if err != nil {
		u.lg.WithError(err).Error("Failed to parse page number")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})

		return
	}
	if pageNumber < 1 {
		u.lg.WithError(err).Error("Failed to parse page number")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})

		return
	}

	userID, ok := ctx.Get(middleware.UserIDKey)
	if !ok {
		u.lg.Error("User ID missing in gin context")
		ctx.JSON(http.StatusInternalServerError, gin.H{})

		return
	}

	newsletters, err := u.getNewslettersByUserID.Handle(ctx, userID.(string), pageSize, pageNumber)
	if err != nil {
		u.lg.WithError(err).Error("Failed to list newsletter")
		ctx.JSON(http.StatusInternalServerError, gin.H{})

		return
	}

	mapped := make([]response.GetNewslettersByUserIDResponse, 0, len(newsletters))
	for _, n := range newsletters {
		mapped = append(mapped, response.CreateGetNewslettersByUserIDResponseFromEntity(n))
	}

	ctx.JSON(http.StatusOK, mapped)
}
