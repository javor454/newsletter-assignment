package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/javor454/newsletter-assignment/app/http_server"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/domain"
	"github.com/javor454/newsletter-assignment/internal/ui/http/request"
	"github.com/javor454/newsletter-assignment/internal/ui/http/response"
)

type GetNewslettersBySubscriberEmailHandler interface {
	Handle(ctx context.Context, email string, pageSize, pageNumber int) ([]*domain.Newsletter, *dto.Pagination, error)
}

type SubscriberController struct {
	lg                                     logger.Logger
	getNewslettersBySubscriberEmailHandler GetNewslettersBySubscriberEmailHandler
}

func NewSubscriberController(
	lg logger.Logger,
	httpServer *http_server.Server,
	gsnbeh GetNewslettersBySubscriberEmailHandler,
) *SubscriberController {
	controller := &SubscriberController{
		getNewslettersBySubscriberEmailHandler: gsnbeh,
		lg:                                     lg,
	}

	httpServer.GetEngine().GET("api/v1/subscribers/:email/newsletters", controller.GetNewslettersBySubscriberEmail)

	return controller
}

//	 TODO dava smysl tady??:

// GetNewslettersBySubscriberEmail
//
//	@Summary	GetNewslettersBySubscriberEmail - retrieve newsletter by subscriber's email
//	@Router		/api/v1/subscribers/{email}/newsletters [get]
//	@Tags		public subscriber
//
//	@Param		Content-Type	header	string	true	"application/json"			default(application/json)
//	@Param		page_size		query	int		true	"Number of items on page"	default(10)	minimum(1)
//	@Param		page_number		query	int		true	"Page number"				default(1)	minimum(1)
//	@Param		email			path	string	true	"Subscribers email"
//
//	@Success	201				"Successfully retrieved newsletters by subscriber email"
//	@Failure	400				{object}	response.Error	"Invalid request with detail"
//	@Failure	500				"Unexpected exception"
func (u *SubscriberController) GetNewslettersBySubscriberEmail(ctx *gin.Context) {
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

	newsletters, pagination, err := u.getNewslettersBySubscriberEmailHandler.Handle(ctx, email, pageSize, pageNumber)
	if err != nil {
		code, body := func(err error) (int, gin.H) {
			return http.StatusInternalServerError, gin.H{}
		}(err)
		u.lg.WithError(err).Error("Failed to get newsletter by subscriber email")
		ctx.JSON(code, body)

		return
	}

	// TODO (nice2have): rm duplicated logic - maybe factory??
	mapped := make([]response.GetNewslettersBySubscriberEmailResponse, 0, len(newsletters))
	for _, n := range newsletters {
		mapped = append(mapped, response.CreateGetNewslettersBySubscriberEmailResponseFromEntity(n))
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
