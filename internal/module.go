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
	"github.com/javor454/newsletter-assignment/internal/infrastructure/auth"
	firebaseinfra "github.com/javor454/newsletter-assignment/internal/infrastructure/firebase"
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
	gnbse := operation.NewGetNewslettersBySubscriberEmail(pgConn)
	guej := operation.NewGetUnsentSubscribedEmailJobsOperation(pgConn)
	uuej := operation.NewUpdateUnsentEmailJobs(pgConn)
	uds := operation.NewUpdateDisableSubscription(pgConn)
	gnbpi := operation.NewGetNewslettersByPublicID(pgConn)

	ms := sendgridinfra.NewMailService(mailClient, appConfig.SendGridTemplateID)

	sc := firebaseinfra.NewSubscriptionCacheManager(fbClient)

	ur := pg.NewUserRepository(cuo, gube)
	nr := pg.NewNewsletterRepository(cno, gnbui, gnbse, gnbpi)
	sr := service.NewSubscriberRepository(lg, pgConn, gnibpi, guej, ms, uuej, appConfig, uds, sc)

	js := auth.NewJwtService(appConfig.JwtSecret)

	hm := healthcheck.NewHealthMonitor(
		healthcheck.NewPgIndicator(pgConn, 5*time.Second),
	)

	ruh := handler.NewRegisterUserHandler(ur, js)
	luh := handler.NewLoginUserHandler(ur, js)
	dth := handler.NewDecodeTokenHandler(js)
	cnh := handler.NewCreateNewsletterHandler(nr)
	gnbuih := handler.NewGetNewslettersByUserIDHandler(nr)
	stnh := handler.NewSubscribeToNewsletterHandler(sr)
	gnbseh := handler.NewGetNewslettersBySubscriberEmailHandler(nr)
	uh := handler.NewUnsubscribeNewsletterHandler(sr)
	pejh := handler.NewProcessEmailJobsHandler(lg, sr)
	pejh.Handle(ctx)
	gnbpih := handler.NewGetNewslettersByPublicIDHandler(nr)

	am := middleware.NewAuthMiddleware(dth, lg)

	controller.NewHealthController(lg, httpServer, hm)
	controller.NewUserController(lg, httpServer, ruh, luh)
	controller.NewNewsletterController(lg, httpServer, cnh, gnbuih, gnbpih, am)
	controller.NewSubscriptionController(lg, httpServer, gnbseh, stnh, uh)
}
