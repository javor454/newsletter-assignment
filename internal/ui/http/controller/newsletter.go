package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/javor454/newsletter-assignment/app/http_server"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/internal/application"
	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/domain"
	"github.com/javor454/newsletter-assignment/internal/ui/http/middleware"
	"github.com/javor454/newsletter-assignment/internal/ui/http/request"
	"github.com/javor454/newsletter-assignment/internal/ui/http/response"
)

type CreateNewsletterHandler interface {
	Handle(ctx context.Context, userID, name string, description *string) error
}

type SubscribeToNewsletterHandler interface {
	Handle(ctx context.Context, newsletterPublicID, email string) error
}

type GetNewslettersByUserIDHandler interface {
	Handle(ctx context.Context, userID string, pageSize, pageNumber int) ([]*domain.Newsletter, *dto.Pagination, error)
}

type NewsletterController struct {
	lg                     logger.Logger
	createNewsletter       CreateNewsletterHandler
	subscribeToNewsletter  SubscribeToNewsletterHandler
	getNewslettersByUserID GetNewslettersByUserIDHandler
}

func NewNewsletterController(
	lg logger.Logger,
	httpServer *http_server.Server,
	cnh CreateNewsletterHandler,
	stn SubscribeToNewsletterHandler,
	gnbui GetNewslettersByUserIDHandler,
	authMiddleware *middleware.AuthMiddleware,
) *NewsletterController {
	controller := &NewsletterController{
		createNewsletter:       cnh,
		subscribeToNewsletter:  stn,
		getNewslettersByUserID: gnbui,
		lg:                     lg,
	}

	httpServer.GetEngine().POST("api/v1/newsletters", authMiddleware.Handle, controller.Create)
	httpServer.GetEngine().POST(
		"api/v1/newsletters/:newsletter_public_id/subscriptions",
		controller.SubscribeToNewsletter,
	)
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
		code, body := func(err error) (int, gin.H) {
			if errors.Is(err, application.UnknownUserError) {
				return http.StatusNotFound, gin.H{"error": "Unknown user"}
			}

			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to create newsletter")
		ctx.JSON(code, body)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func (u *NewsletterController) SubscribeToNewsletter(ctx *gin.Context) {
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

	newsletterID := ctx.Param("newsletter_public_id")
	if newsletterID == "" {
		u.lg.Error("Invalid newsletter_id parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	var req *request.SubscribeToNewsletter
	if err := ctx.ShouldBindJSON(&req); err != nil {
		u.lg.WithError(err).Error("Failed to bind request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	// TODO: handle duplicit subscribe
	if err := u.subscribeToNewsletter.Handle(ctx, newsletterID, req.Email); err != nil {
		code, body := func(err error) (int, gin.H) {
			if errors.Is(err, application.AlreadySubscibedToNewsletterError) {
				return http.StatusConflict, gin.H{"error": "Already subscribed to newsletter"}
			}

			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to subscribe to newsletter")
		ctx.JSON(code, body)

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

	newsletters, pagination, err := u.getNewslettersByUserID.Handle(ctx, userID.(string), pageSize, pageNumber)
	if err != nil {
		u.lg.WithError(err).Error("Failed to list newsletter")
		ctx.JSON(http.StatusInternalServerError, gin.H{})

		return
	}

	mapped := make([]response.GetNewslettersByUserIDResponse, 0, len(newsletters))
	for _, n := range newsletters {
		mapped = append(mapped, response.CreateGetNewslettersByUserIDResponseFromEntity(n))
	}

	ctx.JSON(http.StatusOK, response.PaginatedResponse{
		Data: mapped,
		Pagination: response.Pagination{
			CurrentPage: pagination.CurrentPage,
			PageSize:    pagination.PageSize,
			TotalPages:  pagination.TotalPages,
			TotalItems:  pagination.TotalItems,
			HasPrevious: pagination.HasPrevious,
			HasNext:     pagination.HasNext,
		},
	})
}
