package internal

import (
	"database/sql"
	"time"

	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/healthcheck"
	"github.com/javor454/newsletter-assignment/app/http_server"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/internal/application/handler"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/auth"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/operation"
	"github.com/javor454/newsletter-assignment/internal/ui/http/controller"
	"github.com/javor454/newsletter-assignment/internal/ui/http/middleware"
)

func RegisterDependencies(httpServer *http_server.Server, lg logger.Logger, pgConn *sql.DB, appConfig *config.AppConfig) {
	cuo := operation.NewCreateUser(pgConn)
	gube := operation.NewGetUserByEmail(pgConn)
	cno := operation.NewCreateNewsletter(pgConn)
	gnbui := operation.NewGetNewslettersByUserID(pgConn)
	cs := operation.NewCreateSubscription(pgConn)
	gnibpi := operation.NewGetNewsletterIDByPublicID(pgConn)
	gnbse := operation.NewGetNewslettersBySubscriberEmail(pgConn)

	ur := pg.NewUserRepository(cuo, gube)
	nr := pg.NewNewsletterRepository(cno, gnbui, gnbse)
	sr := pg.NewSubscriberRepository(gnibpi, cs)

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

	am := middleware.NewAuthMiddleware(dth, lg)

	controller.NewHealthController(lg, httpServer, hm)
	controller.NewUserController(lg, httpServer, ruh, luh)
	controller.NewNewsletterController(lg, httpServer, cnh, stnh, gnbuih, am)
	controller.NewSubscriberController(lg, httpServer, gnbseh)
}
