package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/internal/domain"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/operation"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/row"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/sendgrid"
)

type SubscriptionParams struct {
	Email              string `json:"email"`
	NewsletterPublicID string `json:"newsletter_id"`
}

type SubscriptionCache interface {
	CacheSubscription(ctx context.Context, email, newsletterPublicID string) error
}

type SubscriberRepository struct {
	lg                        logger.Logger
	pgConn                    *sql.DB
	getNewsletterByPublicID   *operation.GetNewsletterIDByPublicID
	getUnsendEmailJobs        *operation.GetUnsentSubscribedEmailJobsOperation
	mailService               *sendgrid.MailService
	updateUnsentEmailJobs     *operation.UpdateUnsentEmailJobs
	appConfig                 *config.AppConfig
	updateDisableSubscription *operation.UpdateDisableSubscription
	subscriptionCache         SubscriptionCache
}

func NewSubscriberRepository(
	lg logger.Logger,
	pgConn *sql.DB,
	gn *operation.GetNewsletterIDByPublicID,
	gu *operation.GetUnsentSubscribedEmailJobsOperation,
	ms *sendgrid.MailService,
	uo *operation.UpdateUnsentEmailJobs,
	conf *config.AppConfig,
	uds *operation.UpdateDisableSubscription,
	sc SubscriptionCache,
) *SubscriberRepository {
	lg.Infof("[EMAIL] Sending email: %v", conf.SendMail)
	return &SubscriberRepository{
		lg:                        lg,
		pgConn:                    pgConn,
		getNewsletterByPublicID:   gn,
		getUnsendEmailJobs:        gu,
		mailService:               ms,
		updateUnsentEmailJobs:     uo,
		appConfig:                 conf,
		updateDisableSubscription: uds,
		subscriptionCache:         sc,
	}
}

// TODO (nice2have): simplify somehow

// Subscribe a
func (s *SubscriberRepository) Subscribe(ctx context.Context, subscription *domain.Subscription) error {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond) // TODO: short or long??
	defer cancel()

	idRow, err := s.getNewsletterByPublicID.Execute(
		ctx,
		&operation.GetNewsletterIDByPublicIDParams{PublicID: subscription.NewsletterPublicID().String()},
	)
	if err != nil {
		return err
	}

	// sql.LevelReadCommited
	// - prevents reads from uncommited changes from other txs
	// - allows other transactions to insert jobs simultaneously without blocking each other
	tx, err := s.pgConn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	if err := operation.CreateSubscriptionTx(ctx, tx, &operation.CreateSubscriptionParams{
		ID:              subscription.ID().String(),
		SubscriberEmail: subscription.Email().String(),
		NewsletterID:    idRow.ID,
	}); err != nil {
		return rollback(tx, err)
	}
	paramsJson, err := json.Marshal(SubscriptionParams{
		Email:              subscription.Email().String(),
		NewsletterPublicID: subscription.NewsletterPublicID().String(),
	})
	if err != nil {
		return rollback(tx, err)
	}

	if err := operation.CreateEmailJobTx(ctx, tx, &operation.CreateEmailJobParams{
		ID:     uuid.New().String(),
		Type:   row.SubscriptionType,
		Params: paramsJson,
	}); err != nil {
		return rollback(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit subscription tx: %w", err)
	}

	return nil
}

func (s *SubscriberRepository) ProcessSubscribeEmailJobs(ctx context.Context) error {
	const maxJobs = 100

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	jobs, err := s.getUnsendEmailJobs.Execute(ctx, maxJobs)
	if err != nil {
		return err
	}

	if len(jobs) == 0 {
		return nil
	}
	s.lg.Debugf("[EMAIL] Processing %d jobs...", len(jobs))

	var wg sync.WaitGroup
	var mu sync.Mutex
	processedIDs := make([]string, 0, len(jobs))

	for _, job := range jobs {
		wg.Add(1)

		go func(emailJob *row.EmailJob) {
			defer wg.Done()

			// TODO (nice2have): make ProcessSubscribeEmailJobs generic for all message types
			if job.Type != row.SubscriptionType {
				s.lg.WithField("job_id", job.ID).Errorf("invalid job type on job processing: %s", job.Type)
				return
			}

			var subscribeParams SubscriptionParams
			if err := json.Unmarshal(job.Params, &subscribeParams); err != nil {
				s.lg.WithField("job_id", job.ID).WithError(err).Error("failed to unmarshal subscription job params")
				return
			}

			if s.appConfig.SendMail {
				if err := s.mailService.SendSubscribed(subscribeParams.Email); err != nil {
					s.lg.WithField("job_id", job.ID).WithError(err).Error("failed to send subscribed email")
					return
				}
			}

			mu.Lock()
			processedIDs = append(processedIDs, emailJob.ID)
			mu.Unlock()

			// TODO: if cache fails, processing is still successful to prevent infinite mails
			//  - maybe periodically check data? or introduce integrity hash and invalidate it at the begin of processing
			if err := s.subscriptionCache.CacheSubscription(ctx, subscribeParams.Email, subscribeParams.NewsletterPublicID); err != nil {
				s.lg.WithField("job_id", job.ID).WithError(err).Error("failed to cache subscription")
				return
			}
		}(job)
	}

	wg.Wait()

	if len(processedIDs) == 0 {
		return nil
	}

	ctx, cancel = context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	if err := s.updateUnsentEmailJobs.Execute(ctx, &operation.UpdateUnsentEmailJobsParams{
		JobIDs: processedIDs,
	}); err != nil {
		return err
	}

	return nil
}

func (s *SubscriberRepository) Unsubscribe(ctx context.Context, email *domain.Email, newsletterPublicID *domain.ID) error {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if err := s.updateDisableSubscription.Execute(ctx, &operation.UpdateDisableSubscriptionParams{
		Email:              email.String(),
		NewsletterPublicID: newsletterPublicID.String(),
	}); err != nil {
		return err
	}

	return nil
}

func rollback(tx *sql.Tx, err error) error {
	txErr := tx.Rollback()
	if txErr != nil {
		return fmt.Errorf("failed to rollback transaction: %w", txErr)
	}

	return err
}
