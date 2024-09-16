package internal

import (
	"database/sql"

	"github.com/javor454/newsletter-assignment/app/config"
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

	ur := pg.NewUserRepository(cuo, gube)
	nr := pg.NewNewsletterRepository(cno, gnbui)

	js := auth.NewJwtService(appConfig.JwtSecret)

	ruh := handler.NewRegisterUserHandler(ur, js)
	luh := handler.NewLoginUserHandler(ur, js)
	dth := handler.NewDecodeTokenHandler(js)
	cnh := handler.NewCreateNewsletterHandler(nr)
	lnh := handler.NewGetNewslettersByUserIDHandler(nr)

	am := middleware.NewAuthMiddleware(dth, lg)

	controller.NewUserController(lg, httpServer, ruh, luh)
	controller.NewNewsletterController(lg, httpServer, cnh, lnh, am)
}
