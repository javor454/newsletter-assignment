package internal

import (
	"context"
	"database/sql"
	"time"

	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/firebase"
	"github.com/javor454/newsletter-assignment/app/healthcheck"
	"github.com/javor454/newsletter-assignment/app/http_server"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/app/sendgrid"
	"github.com/javor454/newsletter-assignment/internal/application/handler"
	firebaseinfra "github.com/javor454/newsletter-assignment/internal/infrastructure/firebase"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/jwt"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/operation"
	sendgridinfra "github.com/javor454/newsletter-assignment/internal/infrastructure/sendgrid"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/service"
	"github.com/javor454/newsletter-assignment/internal/ui/http/controller"
	"github.com/javor454/newsletter-assignment/internal/ui/http/middleware"
)

func RegisterDependencies(
	ctx context.Context,
	lg logger.Logger,
	appConfig *config.AppConfig,
	pgConn *sql.DB,
	httpServer *http_server.Server,
	mailClient *sendgrid.Client,
	fbClient *firebase.Client,
) {
	cuo := operation.NewCreateUser(pgConn)
	gube := operation.NewGetUserByEmail(pgConn)
	cno := operation.NewCreateNewsletter(pgConn)
	gnbui := operation.NewGetNewslettersByUserID(pgConn)
	gnibpi := operation.NewGetNewsletterIDByPublicID(pgConn)
	gnbse := operation.NewGetNewslettersBySubscriptionEmail(pgConn)
	guej := operation.NewGetUnsentSubscribedEmailJobsOperation(pgConn)
	uuej := operation.NewUpdateUnsentEmailJobs(pgConn)
	uds := operation.NewUpdateDisableSubscription(pgConn)
	gnbpi := operation.NewGetNewslettersByPublicID(pgConn)

	ms := sendgridinfra.NewMailService(lg, appConfig, mailClient)

	sc := firebaseinfra.NewSubscriptionCacheManager(fbClient)

	ur := pg.NewUserRepository(cuo, gube)
	nr := pg.NewNewsletterRepository(cno, gnbui, gnbse, gnbpi)
	sr := service.NewSubscriberRepository(lg, pgConn, gnibpi, guej, ms, uuej, appConfig, uds, sc)

	tm := jwt.NewTokenManager(appConfig.JwtSecret, appConfig.Host)

	hm := healthcheck.NewHealthMonitor(
		healthcheck.NewPgIndicator(pgConn, 5*time.Second),
	)

	ruh := handler.NewRegisterUserHandler(ur, tm)
	luh := handler.NewLoginUserHandler(ur, tm)
	dth := handler.NewDecodeTokenHandler(tm)
	cnh := handler.NewCreateNewsletterHandler(nr)
	gnbuih := handler.NewGetNewslettersByUserIDHandler(nr)
	stnh := handler.NewSubscribeToNewsletterHandler(tm, sr)
	gnbseh := handler.NewGetNewslettersBySubscriptionEmailHandler(nr)
	uh := handler.NewUnsubscribeNewsletterHandler(sr, tm, sc)
	pejh := handler.NewProcessEmailJobsHandler(lg, sr)
	pejh.Handle(ctx)
	gnbpih := handler.NewGetNewslettersByPublicIDHandler(nr)

	am := middleware.NewAuthMiddleware(dth, lg)

	controller.NewHealthController(lg, httpServer, hm)
	controller.NewUserController(lg, httpServer, ruh, luh)
	controller.NewNewsletterController(lg, httpServer, cnh, gnbuih, gnbpih, am)
	controller.NewSubscriptionController(lg, httpServer, gnbseh, stnh, uh)
}
