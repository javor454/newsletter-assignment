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
	"github.com/javor454/newsletter-assignment/internal/ui/http/request"
	"github.com/javor454/newsletter-assignment/internal/ui/http/response"
)

type GetNewslettersBySubscriptionEmailHandler interface {
	Handle(ctx context.Context, email string, pageSize, pageNumber int) ([]*domain.Newsletter, *dto.Pagination, error)
}

type SubscribeToNewsletterHandler interface {
	Handle(ctx context.Context, newsletterPublicID, email string) error
}

type UnsubscribeNewsletterHandler interface {
	Handle(ctx context.Context, newsletterPublicID, email string) error
}

type SubscriptionController struct {
	lg                                       logger.Logger
	getNewslettersBySubscriptionEmailHandler GetNewslettersBySubscriptionEmailHandler
	subscribeToNewsletter                    SubscribeToNewsletterHandler
	unsubscribeNewsletterHandler             UnsubscribeNewsletterHandler
}

func NewSubscriptionController(
	lg logger.Logger,
	httpServer *http_server.Server,
	gsnbeh GetNewslettersBySubscriptionEmailHandler,
	stnh SubscribeToNewsletterHandler,
	unh UnsubscribeNewsletterHandler,
) *SubscriptionController {
	controller := &SubscriptionController{
		getNewslettersBySubscriptionEmailHandler: gsnbeh,
		lg:                                       lg,
		unsubscribeNewsletterHandler:             unh,
		subscribeToNewsletter:                    stnh,
	}

	httpServer.GetEngine().GET("api/v1/subscriptions/:email/newsletters", controller.GetNewslettersBySubscriptionEmail)
	httpServer.GetEngine().POST(
		"api/v1/newsletters/:newsletter_public_id/subscriptions",
		controller.SubscribeToNewsletter,
	)
	httpServer.GetEngine().DELETE(
		"api/v1/newsletters/:newsletter_public_id/subscriptions/:email",
		controller.UnsubscribeNewsletter,
	)

	return controller
}

// GetNewslettersBySubscriptionEmail
//
//	@Summary	GetNewslettersBySubscriptionEmail - retrieve newsletter by subscriber's email
//	@Router		/api/v1/subscriptions/{email}/newsletters [get]
//	@Tags		public subscription
//	@Produce	json
//
//	@Param		Content-Type	header	string	true	"application/json"			default(application/json)
//	@Param		page_size		query	int		true	"Number of items on page"	default(10)	minimum(1)
//	@Param		page_number		query	int		true	"Page number"				default(1)	minimum(1)
//	@Param		email			path	string	true	"Subscribers email"			default(test@test.com)
//
//	@Success	200				"Successfully retrieved newsletters by subscriber email"
//	@Failure	400				{object}	response.Error	"Invalid request with detail"
//	@Failure	500				"Unexpected exception"
func (u *SubscriptionController) GetNewslettersBySubscriptionEmail(ctx *gin.Context) {
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

	email := ctx.Param("email")
	if email == "" {
		u.lg.Error("Invalid email parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	newsletters, pagination, err := u.getNewslettersBySubscriptionEmailHandler.Handle(ctx, email, pageSize, pageNumber)
	if err != nil {
		code, body := func(err error) (int, gin.H) {
			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to get newsletter by subscriber email")
		ctx.JSON(code, body)

		return
	}

	// TODO (nice2have): rm duplicated logic - maybe factory??
	mapped := make([]*response.PublicNewsletter, 0, len(newsletters))
	for _, n := range newsletters {
		mapped = append(mapped, response.CreatePublicNewsletterResponseFromEntity(n))
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

// SubscribeToNewsletter
//
//	@Summary	SubscribeToNewsletter - used to subscribe to newsletter by email
//	@Router		/api/v1/newsletters/{newsletter_public_id}/subscriptions [post]
//	@Tags		public subscription
//	@Produce	json
//
//	@Param		newsletter_public_id	path	string							true	"Public newsletter identifier"
//	@Param		email					body	request.SubscribeToNewsletter	true	"Subscriber email address"
//
//	@Success	201						"Successfully subscribed to newsletter"
//	@Failure	400						{object}	response.Error	"Invalid request with detail"
//	@Failure	404						{object}	response.Error	"Newsletter not found"
//	@Failure	409						{object}	response.Error	"Already subscribed to newsletter"
//	@Failure	500						"Unexpected exception"
func (u *SubscriptionController) SubscribeToNewsletter(ctx *gin.Context) {
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

// UnsubscribeNewsletter
//
//	@Summary	UnsubscribeNewsletter - used to unsubscribe from newsletter by email
//	@Router		/api/v1/newsletters/{newsletter_public_id}/subscriptions/{email} [delete]
//	@Tags		public subscription
//	@Produce	json
//
//	@Param		Content-Type	header	string	true	"application/json"			default(application/json)
//	@Param		newsletter_public_id	path	string	true	"Public newsletter identifier"
//	@Param		email					path	string	true	"Subscriber email address" default(test@test.com)
//
//	@Success	201						"Successfully unsubscribed from newsletter"
//	@Failure	400						{object}	response.Error	"Invalid request with detail"
//	@Failure	500						"Unexpected exception"
func (u *SubscriptionController) UnsubscribeNewsletter(ctx *gin.Context) {
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

	email := ctx.Param("email")
	if email == "" {
		u.lg.Error("Invalid email parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	if err := u.unsubscribeNewsletterHandler.Handle(ctx, newsletterID, email); err != nil {
		code, body := func(err error) (int, gin.H) {
			if errors.Is(err, application.InvalidUUIDError) {
				return http.StatusBadRequest, gin.H{"error": "Invalid UUID"} // TODO (nice2have): improve message
			}
			if errors.Is(err, application.NewsletterNotFoundError) {
				return http.StatusNotFound, gin.H{"error": "Newsletter not found"}
			}

			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to unsubscribe from newsletter")
		ctx.JSON(code, body)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}
