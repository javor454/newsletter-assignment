package controller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	controllertest "github.com/javor454/newsletter-assignment/test/func/controller"
	"github.com/javor454/newsletter-assignment/test/helper"
	"github.com/stretchr/testify/suite"
)

type NewsletterTestSuite struct {
	suite.Suite
	lg            logger.Logger
	appConf       *config.AppConfig
	pgConn        *sql.DB
	c             *controller.NewsletterController
	am            *middleware.AuthMiddleware
	userIDs       []string
	newsletterIDs []string
}

type newsletterRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type getNewsletterByUserIDResponse struct {
	ID          string    `json:"id"`
	PublicID    string    `json:"public_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type getNewsletterByPublicIDResponse struct {
	PublicID    string    `json:"public_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *NewsletterTestSuite) SetupSuite() {
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

	cn := operation.NewCreateNewsletter(pgConn)
	gn := operation.NewGetNewslettersByUserID(pgConn)
	gns := operation.NewGetNewslettersBySubscriptionEmail(pgConn)
	gnbp := operation.NewGetNewslettersByPublicID(pgConn)

	gnbpi := pg.NewNewsletterRepository(cn, gn, gns, gnbp)

	tm := jwt.NewTokenManager(s.appConf.JwtSecret, s.appConf.Host)

	dth := handler.NewDecodeTokenHandler(tm)
	cnh := handler.NewCreateNewsletterHandler(gnbpi)
	gnbuih := handler.NewGetNewslettersByUserIDHandler(gnbpi)
	gnbpih := handler.NewGetNewsletterByPublicIDHandler(gnbpi)

	s.am = middleware.NewAuthMiddleware(dth, s.lg)

	s.c = controller.NewNewsletterController(s.lg, cnh, gnbuih, gnbpih)
	s.userIDs = make([]string, 0, 2)
	s.newsletterIDs = make([]string, 0, 10)
}

func (s *NewsletterTestSuite) Test_CreateNewsletter_Success() {
	const (
		email                 = "test3@test.com"
		password              = "P@$$w0rD"
		uri                   = "/api/v1/newsletters"
		newsletterName        = "success newsletter 1"
		newsletterDescription = "description 1"
	)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	body := newsletterRequest{
		Name:        newsletterName,
		Description: newsletterDescription,
	}
	jsonBody, err := json.Marshal(&body)
	if err != nil {
		s.T().Fatalf("error marshalling body: %s", err.Error())
	}

	r, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonBody))
	if err != nil {
		s.T().Fatalf("error creating request: %s", err.Error())
	}

	userID := uuid.New().String()
	hash, err := helper.Encrypt(password)
	if err != nil {
		s.T().Fatal(err)
	}
	if err := helper.CreateUser(userID, email, hash, s.pgConn); err != nil {
		s.T().Fatal(err)
	}
	s.userIDs = append(s.userIDs, userID)

	token, err := helper.GenerateJWT(userID, s.appConf.JwtSecret, 5*time.Minute)
	if err != nil {
		s.T().Fatal(err)
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	ctx, engine := gin.CreateTestContext(w)

	ctx.Request = r

	beforeCreate := time.Now()
	engine.Handle(
		http.MethodPost,
		uri,
		s.am.Handle,
		middleware.LoggingMiddleware(s.lg, []string{}),
		s.c.Create,
	)
	engine.HandleContext(ctx)
	afterCreate := time.Now()

	res := w.Result()

	if res.StatusCode != http.StatusCreated {
		s.T().Fatalf("invalid status code: %d", res.StatusCode)
	}

	newsletterRow, err := helper.GetNewslettersByUserID(userID, s.pgConn)
	if err != nil {
		s.T().Fatal(err.Error())
	}
	if len(newsletterRow) != 1 {
		s.T().Fatal("invalid number of saved newsletters")
	}
	s.newsletterIDs = append(s.newsletterIDs, newsletterRow[0].ID)

	s.Equal(newsletterName, newsletterRow[0].Name, "newsletter name mismatch")
	s.Equal(newsletterDescription, newsletterRow[0].Description, "newsletter description mismatch")
	s.True(newsletterRow[0].CreatedAt.After(beforeCreate) && newsletterRow[0].CreatedAt.Before(afterCreate), "invalid creation time")
}

func (s *NewsletterTestSuite) Test_GetNewsletterByUserID_Success() {
	const (
		email                 = "test4@test.com"
		password              = "P@$$w0rD"
		uri                   = "/api/v1/newsletters"
		newsletterName        = "success newsletter 2"
		newsletterDescription = "description 2"
	)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	queryParams := url.Values{}
	queryParams.Add("page_number", "1")
	queryParams.Add("page_size", "10")

	fullURL := fmt.Sprintf("%s?%s", uri, queryParams.Encode())

	r, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		s.T().Fatalf("error creating request: %s", err.Error())
	}

	userID := uuid.New().String()
	hash, err := helper.Encrypt(password)
	if err != nil {
		s.T().Fatalf("encrypt error %s", err.Error())
	}
	if err := helper.CreateUser(userID, email, hash, s.pgConn); err != nil {
		s.T().Fatalf("create user error %s", err.Error())
	}
	s.userIDs = append(s.userIDs, userID)

	newsletterID := uuid.New().String()
	newsletterPublicID := uuid.New().String()

	beforeCreate := time.Now()
	if err := helper.CreateNewsletter(
		newsletterID,
		newsletterPublicID,
		userID,
		newsletterName,
		newsletterDescription,
		s.pgConn,
	); err != nil {
		s.T().Fatalf("creating newsletter error %s", err.Error())
	}
	afterCreate := time.Now()

	s.newsletterIDs = append(s.newsletterIDs, newsletterID)

	token, err := helper.GenerateJWT(userID, s.appConf.JwtSecret, 5*time.Minute)
	if err != nil {
		s.T().Fatalf("generating jwt error %s", err.Error())
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	ctx, engine := gin.CreateTestContext(w)

	ctx.Request = r

	engine.Handle(
		http.MethodGet,
		uri,
		s.am.Handle,
		middleware.LoggingMiddleware(s.lg, []string{}),
		s.c.GetNewslettersByUserID,
	)
	engine.HandleContext(ctx)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		s.T().Fatalf("invalid status code: %d", res.StatusCode)
	}

	var body controllertest.PaginatedResponse[[]getNewsletterByUserIDResponse]
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		s.T().Fatalf("reading body error %s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.T().Fatalf("closing body error %s", err.Error())
		}
	}(res.Body)

	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		s.T().Fatalf("error unmarshalling body: %s", err.Error())
	}

	if len(body.Data) != 1 {
		s.T().Fatalf("invalid number of saved newsletters")
	}

	s.Equal(newsletterID, body.Data[0].ID, "newsletter id mismatch")
	s.Equal(newsletterPublicID, body.Data[0].PublicID, "newsletter public id mismatch")
	s.Equal(newsletterName, body.Data[0].Name, "newsletter name mismatch")
	s.Equal(newsletterDescription, *body.Data[0].Description, "newsletter description mismatch")
	s.True(body.Data[0].CreatedAt.After(beforeCreate) && body.Data[0].CreatedAt.Before(afterCreate), "invalid creation time")
}

func (s *NewsletterTestSuite) Test_GetNewsletterByPublicID_Success() {
	const (
		email                 = "test5@test.com"
		password              = "P@$$w0rD"
		uri                   = "/api/v1/newsletters"
		newsletterName        = "success newsletter 3"
		newsletterDescription = "description 3"
	)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	queryParams := url.Values{}
	queryParams.Add("page_number", "1")
	queryParams.Add("page_size", "10")

	userID := uuid.New().String()
	hash, err := helper.Encrypt(password)
	if err != nil {
		s.T().Fatalf("encrypt error %s", err.Error())
	}
	if err := helper.CreateUser(userID, email, hash, s.pgConn); err != nil {
		s.T().Fatalf("create user error %s", err.Error())
	}
	s.userIDs = append(s.userIDs, userID)

	newsletterID := uuid.New().String()
	newsletterPublicID := uuid.New().String()

	beforeCreate := time.Now()
	if err := helper.CreateNewsletter(
		newsletterID,
		newsletterPublicID,
		userID,
		newsletterName,
		newsletterDescription,
		s.pgConn,
	); err != nil {
		s.T().Fatalf("creating newsletter error %s", err.Error())
	}
	afterCreate := time.Now()

	s.newsletterIDs = append(s.newsletterIDs, newsletterID)

	fullURL := fmt.Sprintf("%s/%s?%s", uri, newsletterPublicID, queryParams.Encode())

	r, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		s.T().Fatalf("error creating request: %s", err.Error())
	}

	token, err := helper.GenerateJWT(userID, s.appConf.JwtSecret, 5*time.Minute)
	if err != nil {
		s.T().Fatalf("generating jwt error %s", err.Error())
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	ctx, engine := gin.CreateTestContext(w)

	ctx.Request = r

	engine.Handle(
		http.MethodGet,
		fmt.Sprintf("%s/%s", uri, ":public_id"),
		s.am.Handle,
		middleware.LoggingMiddleware(s.lg, []string{}),
		s.c.GetNewsletterByPublicID,
	)
	engine.HandleContext(ctx)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		s.T().Fatalf("invalid status code: %d", res.StatusCode)
	}

	var body getNewsletterByPublicIDResponse
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		s.T().Fatalf("reading body error %s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.T().Fatalf("closing body error %s", err.Error())
		}
	}(res.Body)

	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		s.T().Fatalf("error unmarshalling body: %s", err.Error())
	}

	s.Equal(newsletterPublicID, body.PublicID, "newsletter public id mismatch")
	s.Equal(newsletterName, body.Name, "newsletter name mismatch")
	s.Equal(newsletterDescription, *body.Description, "newsletter description mismatch")
	s.True(body.CreatedAt.After(beforeCreate) && body.CreatedAt.Before(afterCreate), "invalid creation time")
}

func (s *NewsletterTestSuite) TearDownSuite() {
	if err := helper.RemoveNewsletterByID(s.newsletterIDs, s.pgConn); err != nil {
		s.T().Fatal(err)
	}
	if err := helper.RemoveUsersByUserID(s.userIDs, s.pgConn); err != nil {
		s.T().Fatal(err)
	}
	if err := s.pgConn.Close(); err != nil {
		s.T().Fatalf("pgConn close failed: %s", err.Error())
	}
}

func TestNewsletterSuite(t *testing.T) {
	suite.Run(t, new(NewsletterTestSuite))
}
