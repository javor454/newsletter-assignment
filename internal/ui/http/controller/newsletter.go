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

// Create
//
//	@Summary	Create - used to create new newsletter
//	@Router		/api/v1/newsletters [post]
//	@Tags		newsletter
//
//	@Param		Content-Type	header	string							true	"application/json"	default(application/json)
//	@Param		Authorization	header	string							true	"Bearer <token>"	default(Bearer )
//	@Param		Newsletter		body	request.CreateNewsletterRequest	true	"Newsletter data to create"
//
//	@Success	201				"Newsletter was successfully created"
//	@Failure	400				{object}	response.Error	"Invalid request with detail"
//	@Failure	404				{object}	response.Error	"Unknown user"
//	@Failure	500				"Unexpected exception"
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
			if errors.Is(err, application.InvalidUUIDError) {
				return http.StatusBadRequest, gin.H{"error": err.Error()}
			}
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

// SubscribeToNewsletter
//
//	@Summary	Create - used to create new newsletter
//	@Router		/api/v1/newsletters/:newsletter_public_id/subscriptions [post]
//	@Tags		newsletter
//
//	@Param		Content-Type			header	string							true	"application/json"	default(application/json)
//	@Param		Authorization			header	string							true	"Bearer <token>"	default(Bearer )
//	@Param		newsletter_public_id	path	string							true	"Public newsletter identifier"
//	@Param		email					body	request.SubscribeToNewsletter	true	"Subscriber email address"
//
//	@Success	201						"Successfully subscribed to newsletter"
//	@Failure	400						{object}	response.Error	"Invalid request with detail"
//	@Failure	404						{object}	response.Error	"Newsletter not found"
//	@Failure	409						{object}	response.Error	"Already subscribed to newsletter"
//	@Failure	500						"Unexpected exception"
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

	if err := u.subscribeToNewsletter.Handle(ctx, newsletterID, req.Email); err != nil {
		code, body := func(err error) (int, gin.H) {
			if errors.Is(err, application.InvalidUUIDError) {
				return http.StatusBadRequest, gin.H{"error": err.Error()}
			}
			if errors.Is(err, application.AlreadySubscibedToNewsletterError) {
				return http.StatusConflict, gin.H{"error": "Already subscribed to newsletter"}
			}
			if errors.Is(err, application.NewsletterNotFoundError) {
				return http.StatusNotFound, gin.H{"error": "Newsletter not found"}
			}

			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to subscribe to newsletter")
		ctx.JSON(code, body)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

// GetNewslettersByUserID
//
//	@Summary	GetNewslettersByUserID - retrieve newsletter by creator's user ID
//	@Router		/api/v1/newsletters [get]
//	@Tags		newsletter
//
//	@Param		Content-Type	header	string	true	"application/json"			default(application/json)
//	@Param		Authorization	header	string	true	"Bearer <token>"			default(Bearer )
//	@Param		page_size		query	int		true	"Number of items on page"	default(10)	minimum(1)
//	@Param		page_number		query	int		true	"Page number"				default(1)	minimum(1)
//
//	@Success	201				"Successfully retrieved newsletters by user ID"
//	@Failure	400				{object}	response.Error	"Invalid request with detail"
//	@Failure	409				{object}	response.Error	"Already subscribed to newsletter"
//	@Failure	500				"Unexpected exception"
func (u *NewsletterController) GetNewslettersByUserID(ctx *gin.Context) {
	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil {
		u.lg.WithError(err).Error("Failed to parse page size")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid pageSize"})

		return
	}
	// TODO:a
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
		code, body := func(err error) (int, gin.H) {
			if errors.Is(err, application.InvalidUUIDError) {
				return http.StatusBadRequest, gin.H{"error": err.Error()}
			}

			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to list newsletter")
		ctx.JSON(code, body)

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
