package handler

import (
	"context"
	"time"

	"github.com/javor454/newsletter-assignment/app/logger"
)

type ProcessEmailJobsService interface {
	ProcessSubscribeEmailJobs(ctx context.Context) error
}

// ProcessEmailJobsHandler processes new email jobs repeatedly until context done is signalled
type ProcessEmailJobsHandler struct {
	lg               logger.Logger
	processEmailJobs ProcessEmailJobsService
}

func NewProcessEmailJobsHandler(lg logger.Logger, processEmailJobs ProcessEmailJobsService) *ProcessEmailJobsHandler {
	return &ProcessEmailJobsHandler{
		lg:               lg,
		processEmailJobs: processEmailJobs,
	}
}

func (h *ProcessEmailJobsHandler) Handle(ctx context.Context) {
	go func() {
		h.lg.Info("[EMAIL] Starting email job processing...")
		for {
			select {
			case <-ctx.Done():
				h.lg.Debug("[EMAIL] Processing stopped")
				return
			case <-time.After(1 * time.Minute):
				h.lg.Debug("[EMAIL] Processing new job batch...")
				if err := h.processEmailJobs.ProcessSubscribeEmailJobs(ctx); err != nil {
					// TODO: if one job fails it keeps running infinitely - can happen only in case of INTERNAL
					h.lg.WithError(err).Error("[EMAIL] Error processing batch")
				}
			}
		}

	}()
}
