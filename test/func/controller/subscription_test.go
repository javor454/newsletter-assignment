package controller_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/firebase"
	"github.com/javor454/newsletter-assignment/app/logger"
	pgapp "github.com/javor454/newsletter-assignment/app/pg"
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
	"github.com/javor454/newsletter-assignment/test/helper"
	"github.com/stretchr/testify/suite"
)

type SubscriptionTestSuite struct {
	suite.Suite
	lg              logger.Logger
	appConf         *config.AppConfig
	pgConn          *sql.DB
	c               *controller.SubscriptionController
	am              *middleware.AuthMiddleware
	userIDs         []string
	newsletterIDs   []string
	subscriptionIDs []string
}

type subscribeRequest struct {
	Email string `json:"email"`
}

func (s *SubscriptionTestSuite) SetupSuite() {
	ctx := context.Background()
	s.appConf = helper.NewAppConfig()
	fbConfig := helper.NewFirebaseConfig()
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
	fbClient, err := firebase.NewClient(s.lg, ctx, fbConfig)
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	mailClient := sendgrid.NewClient(s.lg, s.appConf)

	cn := operation.NewCreateNewsletter(pgConn)
	gn := operation.NewGetNewslettersByUserID(pgConn)
	gns := operation.NewGetNewslettersBySubscriptionEmail(pgConn)
	gnbp := operation.NewGetNewslettersByPublicID(pgConn)
	gnibp := operation.NewGetNewsletterIDByPublicID(pgConn)
	gu := operation.NewGetUnsentSubscribedEmailJobsOperation(pgConn)
	uo := operation.NewUpdateUnsentEmailJobs(pgConn)
	uds := operation.NewUpdateDisableSubscription(pgConn)

	sc := firebaseinfra.NewSubscriptionCacheManager(fbClient)
	tm := jwt.NewTokenManager(s.appConf.JwtSecret, s.appConf.Host)

	nr := pg.NewNewsletterRepository(cn, gn, gns, gnbp)
	ms := sendgridinfra.NewMailService(s.lg, s.appConf, mailClient)
	sr := service.NewSubscriberRepository(s.lg, s.pgConn, gnibp, gu, ms, uo, s.appConf, uds, sc)

	dth := handler.NewDecodeTokenHandler(tm)
	unh := handler.NewUnsubscribeNewsletterHandler(sr, tm, sc)
	stnh := handler.NewSubscribeToNewsletterHandler(tm, sr)
	gsnbeh := handler.NewGetNewslettersBySubscriptionEmailHandler(nr)

	s.am = middleware.NewAuthMiddleware(dth, s.lg)

	s.c = controller.NewSubscriptionController(s.lg, gsnbeh, stnh, unh)
	s.userIDs = make([]string, 0, 2)
	s.newsletterIDs = make([]string, 0, 10)
	s.subscriptionIDs = make([]string, 0, 10)
}

func (s *SubscriptionTestSuite) Test_SubscribeToNewsletter_Success() {
	const (
		email                 = "test6@test.com"
		subscriberEmail       = "subscriber1@test.com"
		password              = "P@$$w0rD"
		newsletterName        = "success newsletter 4"
		newsletterDescription = "description 4"
	)

	// fixtures
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

	s.newsletterIDs = append(s.newsletterIDs, newsletterID)

	// setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	body := subscribeRequest{
		Email: subscriberEmail,
	}
	jsonBody, err := json.Marshal(&body)
	if err != nil {
		s.T().Fatalf("error marshalling body: %s", err.Error())
	}

	r, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/api/v1/newsletters/%s/subscriptions", newsletterPublicID),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		s.T().Fatalf("error creating request: %s", err.Error())
	}

	r.Header.Set("Content-Type", "application/json")

	ctx, engine := gin.CreateTestContext(w)

	ctx.Request = r

	beforeCreate := time.Now()
	engine.Handle(
		http.MethodPost,
		"/api/v1/newsletters/:newsletter_public_id/subscriptions",
		middleware.LoggingMiddleware(s.lg, []string{}),
		s.c.SubscribeToNewsletter,
	)
	engine.HandleContext(ctx)
	afterCreate := time.Now()

	res := w.Result()

	if res.StatusCode != http.StatusCreated {
		s.T().Fatalf("invalid status code: %d", res.StatusCode)
	}

	subscriptionRows, err := helper.GetSubscriptionByNewsletterID(newsletterID, s.pgConn)
	if err != nil {
		s.T().Fatal(err.Error())
	}
	if len(subscriptionRows) != 1 {
		s.T().Fatal("invalid number of saved subscriptions")
	}
	s.subscriptionIDs = append(s.subscriptionIDs, subscriptionRows[0].ID)

	fmt.Println(beforeCreate, afterCreate)
	s.Equal(subscriberEmail, subscriptionRows[0].SubscriberEmail)
	s.Equal(newsletterID, subscriptionRows[0].NewsletterID)
	s.True(subscriptionRows[0].CreatedAt.After(beforeCreate) && subscriptionRows[0].CreatedAt.Before(afterCreate), "invalid creation time")
	s.Nil(subscriptionRows[0].DisabledAt)

	subject, err := helper.ParseJWT(s.appConf.JwtSecret, subscriptionRows[0].Token)
	if err != nil {
		s.T().Fatalf("parse jwt error: %s", err.Error())
	}
	s.Equal(subscriberEmail, subject)
}

// func (s *SubscriptionTestSuite) Test_GetNewsletterByUserID_Success() {
// 	const (
// 		email                 = "test4@test.com"
// 		password              = "P@$$w0rD"
// 		uri                   = "/api/v1/newsletters"
// 		newsletterName        = "success newsletter 2"
// 		newsletterDescription = "description 2"
// 	)
//
// 	gin.SetMode(gin.TestMode)
// 	w := httptest.NewRecorder()
//
// 	queryParams := url.Values{}
// 	queryParams.Add("page_number", "1")
// 	queryParams.Add("page_size", "10")
//
// 	fullURL := fmt.Sprintf("%s?%s", uri, queryParams.Encode())
//
// 	r, err := http.NewRequest(http.MethodGet, fullURL, nil)
// 	if err != nil {
// 		s.T().Fatalf("error creating request: %s", err.Error())
// 	}
//
// 	userID := uuid.New().String()
// 	hash, err := helper.Encrypt(password)
// 	if err != nil {
// 		s.T().Fatalf("encrypt error %s", err.Error())
// 	}
// 	if err := helper.CreateUser(userID, email, hash, s.pgConn); err != nil {
// 		s.T().Fatalf("create user error %s", err.Error())
// 	}
// 	s.userIDs = append(s.userIDs, userID)
//
// 	newsletterID := uuid.New().String()
// 	newsletterPublicID := uuid.New().String()
//
// 	beforeCreate := time.Now()
// 	if err := helper.CreateNewsletter(
// 		newsletterID,
// 		newsletterPublicID,
// 		userID,
// 		newsletterName,
// 		newsletterDescription,
// 		s.pgConn,
// 	); err != nil {
// 		s.T().Fatalf("creating newsletter error %s", err.Error())
// 	}
// 	afterCreate := time.Now()
//
// 	s.newsletterIDs = append(s.newsletterIDs, newsletterID)
//
// 	token, err := helper.GenerateJWT(userID, s.appConf.JwtSecret, 5*time.Minute)
// 	if err != nil {
// 		s.T().Fatalf("generating jwt error %s", err.Error())
// 	}
//
// 	r.Header.Set("Content-Type", "application/json")
// 	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
//
// 	ctx, engine := gin.CreateTestContext(w)
//
// 	ctx.Request = r
//
// 	engine.Handle(
// 		http.MethodGet,
// 		uri,
// 		s.am.Handle,
// 		middleware.LoggingMiddleware(s.lg, []string{}),
// 		s.c.GetNewslettersByUserID,
// 	)
// 	engine.HandleContext(ctx)
//
// 	res := w.Result()
//
// 	if res.StatusCode != http.StatusOK {
// 		s.T().Fatalf("invalid status code: %d", res.StatusCode)
// 	}
//
// 	var body controllertest.PaginatedResponse[[]getNewsletterByUserIDResponse]
// 	bodyBytes, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		s.T().Fatalf("reading body error %s", err.Error())
// 	}
// 	defer func(Body io.ReadCloser) {
// 		err := Body.Close()
// 		if err != nil {
// 			s.T().Fatalf("closing body error %s", err.Error())
// 		}
// 	}(res.Body)
//
// 	if err := json.Unmarshal(bodyBytes, &body); err != nil {
// 		s.T().Fatalf("error unmarshalling body: %s", err.Error())
// 	}
//
// 	if len(body.Data) != 1 {
// 		s.T().Fatalf("invalid number of saved newsletters")
// 	}
//
// 	s.Equal(newsletterID, body.Data[0].ID, "newsletter id mismatch")
// 	s.Equal(newsletterPublicID, body.Data[0].PublicID, "newsletter public id mismatch")
// 	s.Equal(newsletterName, body.Data[0].Name, "newsletter name mismatch")
// 	s.Equal(newsletterDescription, *body.Data[0].Description, "newsletter description mismatch")
// 	s.True(body.Data[0].CreatedAt.After(beforeCreate) && body.Data[0].CreatedAt.Before(afterCreate), "invalid creation time")
// }

func (s *SubscriptionTestSuite) TearDownSuite() {
	if err := helper.RemoveSubscriptionsByID(s.subscriptionIDs, s.pgConn); err != nil {
		s.T().Fatal(err)
	}
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

func TestSubscriptionSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionTestSuite))
}
