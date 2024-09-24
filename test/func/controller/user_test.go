package controller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/logger"
	pgapp "github.com/javor454/newsletter-assignment/app/pg"
	"github.com/javor454/newsletter-assignment/internal/application/handler"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/jwt"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/operation"
	"github.com/javor454/newsletter-assignment/internal/ui/http/controller"
	"github.com/javor454/newsletter-assignment/internal/ui/http/middleware"
	"github.com/javor454/newsletter-assignment/test/helper"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	lg      logger.Logger
	appConf *config.AppConfig
	pgConn  *sql.DB
	c       *controller.UserController
	userIDs []string
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *UserTestSuite) SetupSuite() {
	s.appConf = helper.NewAppConfig()
	location, err := time.LoadLocation(s.appConf.Timezone)
	if err != nil {
		panic("failed to load timezone")
	}
	time.Local = location
	pgConfig := helper.NewPostgresConfig()
	s.lg = logger.NewLogger(s.appConf)
	pgConn, err := pgapp.NewConnection(s.lg, pgConfig)
	if err != nil {
		s.lg.WithError(err).Fatal("pg connection init failed")
	}
	s.pgConn = pgConn
	if err := pgapp.MigrationsUp(s.lg, pgConfig, pgConn); err != nil {
		s.lg.WithError(err).Fatal("pg migrations failed")
	}

	cuo := operation.NewCreateUser(pgConn)
	gube := operation.NewGetUserByEmail(pgConn)

	ur := pg.NewUserRepository(cuo, gube)
	tm := jwt.NewTokenManager(s.appConf.JwtSecret, s.appConf.Host)

	ruh := handler.NewRegisterUserHandler(ur, tm)
	luh := handler.NewLoginUserHandler(ur, tm)

	s.c = controller.NewUserController(s.lg, ruh, luh)
	s.userIDs = make([]string, 0, 2)
}

func (s *UserTestSuite) Test_RegisterUser_Success() {
	const (
		email    = "test1@test.com"
		password = "P@$$w0rD"
		uri      = "/api/v1/users/register"
	)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	body := userRequest{
		Email:    email,
		Password: password,
	}
	jsonBody, err := json.Marshal(&body)
	if err != nil {
		s.T().Fatalf("error marshalling body: %s", err.Error())
	}

	r, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonBody))
	if err != nil {
		s.T().Fatalf("error creating request: %s", err.Error())
	}

	r.Header.Set("Content-Type", "application/json")

	ctx, engine := gin.CreateTestContext(w)

	ctx.Request = r

	beforeRegister := time.Now()
	engine.Handle(
		http.MethodPost,
		uri,
		middleware.LoggingMiddleware(s.lg, []string{}),
		s.c.Register,
	)
	engine.HandleContext(ctx)
	afterRegister := time.Now()

	res := w.Result()

	if res.StatusCode != http.StatusCreated {
		s.T().Fatalf("invalid status code: %d", res.StatusCode)
	}

	userRow, err := helper.GetUserByEmail(email, s.pgConn)
	if err != nil {
		s.T().Fatal(err.Error())
	}
	s.userIDs = append(s.userIDs, userRow.ID)

	token := res.Header.Get("Authorization")
	parts := strings.Split(token, " ")

	userID, err := helper.ParseJWT(s.appConf.JwtSecret, parts[1])
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(userRow.ID, userID, "invalid userID")
	s.True(
		userRow.CreatedAt.After(beforeRegister) && userRow.CreatedAt.Before(afterRegister),
		"invalid creation time",
	)
	s.True(helper.IsEqual(userRow.PasswordHash, password), "invalid match between hash and password")
}

func (s *UserTestSuite) Test_LoginUser_Success() {
	const (
		email    = "test2@test.com"
		password = "P@$$w0rD"
		uri      = "/api/v1/users/login"
	)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	body := userRequest{
		Email:    email,
		Password: password,
	}
	jsonBody, err := json.Marshal(&body)
	if err != nil {
		s.T().Fatalf("error marshalling body: %s", err.Error())
	}

	r, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonBody))
	if err != nil {
		s.T().Fatalf("error creating request: %s", err.Error())
	}

	r.Header.Set("Content-Type", "application/json")

	ctx, engine := gin.CreateTestContext(w)

	ctx.Request = r

	hash, err := helper.Encrypt(password)
	if err != nil {
		s.T().Fatal(err)
	}

	userID := uuid.New().String()
	if err := helper.CreateUser(userID, email, hash, s.pgConn); err != nil {
		s.T().Fatal(err)
	}
	s.userIDs = append(s.userIDs, userID)

	engine.Handle(
		http.MethodPost,
		uri,
		middleware.LoggingMiddleware(s.lg, []string{}),
		s.c.Login,
	)
	engine.HandleContext(ctx)

	res := w.Result()

	s.Equal(http.StatusCreated, res.StatusCode, "invalid status code")
}

func (s *UserTestSuite) TearDownSuite() {
	if err := helper.RemoveUsersByUserID(s.userIDs, s.pgConn); err != nil {
		s.T().Fatal(err)
	}
	if err := s.pgConn.Close(); err != nil {
		s.T().Fatalf("pgConn close failed: %s", err.Error())
	}
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
