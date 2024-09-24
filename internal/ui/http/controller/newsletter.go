package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

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

type GetNewslettersByUserIDHandler interface {
	Handle(ctx context.Context, userID string, pageSize, pageNumber int) ([]*domain.Newsletter, *dto.Pagination, error)
}

type GetNewslettersByPublicIDHandler interface {
	Handle(ctx context.Context, publicID string) (*domain.Newsletter, error)
}

type NewsletterController struct {
	lg                       logger.Logger
	createNewsletter         CreateNewsletterHandler
	getNewslettersByUserID   GetNewslettersByUserIDHandler
	getNewslettersByPublicID GetNewslettersByPublicIDHandler
}

func NewNewsletterController(
	lg logger.Logger,
	cnh CreateNewsletterHandler,
	gnbui GetNewslettersByUserIDHandler,
	gnbpih GetNewslettersByPublicIDHandler,
) *NewsletterController {
	controller := &NewsletterController{
		createNewsletter:         cnh,
		getNewslettersByUserID:   gnbui,
		lg:                       lg,
		getNewslettersByPublicID: gnbpih,
	}

	return controller
}

func (u *NewsletterController) RegisterNewsletterController(
	authMiddleware *middleware.AuthMiddleware,
	httpServer *http_server.Server,
) {
	httpServer.GetEngine().POST("api/v1/newsletters", authMiddleware.Handle, u.Create)
	httpServer.GetEngine().GET("api/v1/newsletters", authMiddleware.Handle, u.GetNewslettersByUserID)

	httpServer.GetEngine().GET("api/v1/newsletters/:public_id", u.GetNewsletterByPublicID)

}

// Create
//
//	@Summary	Create new newsletter
//	@Router		/api/v1/newsletters [post]
//	@Tags		newsletter
//	@Accepts	json
//	@Produce	json
//
//	@Param		Authorization	header	string							true	"Bearer <token>"	default(Bearer )
//	@Param		Newsletter		body	request.CreateNewsletterRequest	true	"Newsletter data to create"
//
//	@Success	201				"Newsletter was successfully created"
//	@Failure	400				{object}	response.Error	"Invalid request with detail"
//	@Failure	401				"Unauthorized"
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
				return http.StatusUnauthorized, gin.H{}
			}

			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to create newsletter")
		ctx.JSON(code, body)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

// GetNewslettersByUserID
//
//	@Summary	Retrieve newsletter by creator's user ID
//	@Router		/api/v1/newsletters [get]
//	@Tags		newsletter
//	@Accepts	json
//	@Produce	json
//
//	@Param		Content-Type	header		string						true	"application/json"			default(application/json)
//	@Param		Authorization	header		string						true	"Bearer <token>"			default(Bearer )
//	@Param		page_size		query		int							true	"Number of items on page"	default(10)	minimum(1)
//	@Param		page_number		query		int							true	"Page number"				default(1)	minimum(1)
//
//	@Success	200				{object}	response.InternalNewsletter	"Successfully retrieved newsletters by user ID"
//	@Failure	400				{object}	response.Error				"Invalid request with detail"
//	@Failure	500				"Unexpected exception"
func (u *NewsletterController) GetNewslettersByUserID(ctx *gin.Context) {
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
		code, body := func(err error) (int, gin.H) {
			if errors.Is(err, application.InvalidUUIDError) {
				return http.StatusBadRequest, gin.H{"error": err.Error()}
			}

			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to get newsletters by user ID")
		ctx.JSON(code, body)

		return
	}

	mapped := make([]*response.InternalNewsletter, 0, len(newsletters))
	for _, n := range newsletters {
		mapped = append(mapped, response.CreateInternalNewsletterResponseFromEntity(n))
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

// GetNewsletterByPublicID
//
//	@Summary	Retrieve newsletter by its public ID
//	@Router		/api/v1/newsletters/{public_id} [get]
//	@Tags		public newsletter
//	@Accepts	json
//	@Produce	json
//
//	@Param		Content-Type	header		string						true	"application/json"	default(application/json)
//	@Param		public_id		path		string						true	"Newsletter public ID"
//
//	@Success	200				{object}	response.PublicNewsletter	"Successfully retrieved newsletter by public ID"
//	@Failure	400				{object}	response.Error				"Invalid request with detail"
//	@Failure	500				"Unexpected exception"
func (u *NewsletterController) GetNewsletterByPublicID(ctx *gin.Context) {
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

	publicID := ctx.Param("public_id")
	if publicID == "" {
		u.lg.Error("Invalid public_id parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	newsletter, err := u.getNewslettersByPublicID.Handle(ctx, publicID)
	if err != nil {
		code, body := func(err error) (int, gin.H) {
			if errors.Is(err, application.InvalidUUIDError) {
				return http.StatusBadRequest, gin.H{"error": err.Error()}
			}

			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to get newsletters by public ID")
		ctx.JSON(code, body)

		return
	}

	ctx.JSON(http.StatusOK, response.PublicNewsletter{
		PublicID:    newsletter.PublicID().String(),
		Name:        newsletter.Name(),
		Description: newsletter.Description(),
		CreatedAt:   newsletter.CreatedAt().Format(time.RFC3339),
	})
}
