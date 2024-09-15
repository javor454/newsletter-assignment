package internal

import (
	"database/sql"

	"github.com/javor454/newsletter-assignment/app/http-server"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/internal/application/handler"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/operation"
	"github.com/javor454/newsletter-assignment/internal/ui/http/controller"
)

func RegisterDependencies(httpServer *http_server.HttpServer, lg logger.Logger, pgConn *sql.DB) {
	cuo := operation.NewCreateUser(pgConn)
	us := pg.NewUserService(cuo)
	ruh := handler.NewRegisterUserHandler(us)
	luh := handler.NewLoginUserHandler()
	controller.NewUserController(lg, httpServer, ruh, luh)
}
