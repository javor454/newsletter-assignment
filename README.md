# Newsletter assignment

## Prerequisities
- gnu make
- docker + docker compose

## Setup
- run `make up`

## API documentation
- available on 
  - development: http://localhost:8080/api/docs/index.html
  - production: TODO

## Architecture
- hexagon
- solid
- transaction outbox pattern

## Features
### Users
#### Registration
- HTTP API designed by REST principles
- public endpoint
- POST `api/v1/users/register`
- scenarios
  - success scenario
    - in request send with email and password
    - validate email
    - encrypt password with bcrypt (hash password + salt)
    - save to postgres
    - generate JWT token containing userID, issued at and expiration timestamp
    - respond with token in authorization header
  - fail scenarios
    - in case of invalid request, receive 400
    - for registration with taken email, receive 409 response

#### Login
- HTTP API designed by REST principles
- public endpoint
- POST `api/v1/users/login`
- scenarios
  - success scenario
    - in request send with email and password
    - get password from database by email
    - hash password from request and compare with password in db
    - generate JWT token containing userID, issued at and expiration timestamp
    - receive Bearer token in response header
  - fail scenarios
    - in case of invalid request, respond with 400
    - in case of invalid credentials (non-registered email, invalid email x password match), receive 401 response

### Newsletter
#### Create newsletter
- HTTP API designed by REST principles
- secured endpoint
- POST `api/v1/newsletters`
- scenarios
  - success scenario
    - use Bearer token for auth in Authorization header
    - in request send name and description
    - create unique UUID for newsletter identification
    - save to postgres
  - fail scenarios
    - in case of invalid request, receive 400
    - in case user is not found in db, respond with 401

#### Get newsletter by user
- HTTP API designed by REST principles
- secured endpoint
- paginated
- GET `api/v1/newsletters`
- success scenario
  - use Bearer token for auth in Authorization header
  - retrieve paginated list of newsletters by user id in token
- fail scenario
  - in case of invalid request, receive 400

### Subscriptions
#### Get newsletter by subscriber email
- HTTP API designed by REST principles
- public endpoint
- paginated
- GET `api/v1/subscribers/:email/newsletters`
- success scenario
  - in path parameter send subscriber email
  - retrieve paginated list of newsletters by email
- fail scenario
  - in case of invalid request, receive 400

#### Subscribe to newsletter
- HTTP API designed by REST principles
- public endpoint
- POST `api/v1/newsletters/:newsletter_public_id/subscriptions`
- success scenario
  - in path parameter send newsletter public id

#### Unsubscribe from newsletter
- HTTP API designed by REST principles
- public endpoint
- DELETE `api/v1/newsletters/:newsletter_public_id/subscriptions/:email`
- success scenario
  - in path parameter send newsletter public id and email

## Flows
- registrations
  - register endpoint
- login
  - login endpoint
- create newsletter
  - register / login
  - create newsletter endpoint
- get user newsletters
  - register / login
  - create newsletter endpoint
  - get user newsletters endpoint
- subscribe to newsletter
  - subscribe endpoint
- get subscribed newsletters
  - subscribe endpoint
  - get subscribed newsletters
    - http endpoint to get newsletters by email
    - firebase get public IDS + http endpoint to get newsletters by public ID 
- unsubscribe from newsletter
  - unsubscribe endpoint

## TODOS for PROD
- system tests (func, unit, integration)
  - jwt token
    - expiration
    - malformed
  - all happy paths + application errors
- infrastructure
  - build binary
  - digital ocean VPS
  - DB? either droplet or run in docker compose
  - setup firebase project permissions dev/prod
  - envs
  - secrets
    - firebase service account key
    - sendgrid api key
    - jwt
  - ssl/tls
  - domain
  - tagging
  - zero downtime deploy
  - reverse proxy with rate limiting
  - github actions
    - lint
    - vulnerability check
    - build
    - test
    - deploy
- features
  - unsubscription link
  - unique code to pair together subscription and unsubscription links
## TODOS extras
- basic performance test
- simple backoffice
- replace GET http endpoints with one public and one private graphql
- improve project structure

## Development plan
- Setup project via docker compose ✅
  - go app ✅
  - postgre ✅
  - firebase ✅
    - configure SDK ✅
  - email server ✅
- air hot rebuild in local env ✅
- Setup other dependencies ✅
  - email server ✅
    - SendGrid or AWS SES ✅
- db migrations
  - structures  ✅
  - data - not necessary ✅
- MVP go app
  - REST
    - user registration ✅
      - db storing password ✅ 
    - user authentication JWT ✅
    - user authorization ✅
    - newsletter management ✅
    - subscription management ✅
    - endpoint na subscribe ktery prijme public newsletter id a email subscribera WIP  ✅
    - email functionality ✅
      - queue ✅
      - disable for development env ✅ 
    - api documentation - swagger ✅
      - definition ✅  
      - test if works ✅
    - healthcheck ✅
      - pg ✅
    - graceful shutdown ✅
      - pg ✅
    - logs `✅`
    - panic recovery ✅
    - security
      - rate limiting
      - cors ✅
    - pagination ✅
    - handle errors ✅
      - registration : invalid password or email ✅
      - registration with same email ✅
    - verify timeouts in infra ✅
  - testing
    - system tests (func, unit, integration)
    - jwt token
      - expiration
      - malformed
    - all happy paths + application errors
- FIX
  - check volumes ✅
- deployment
  - digital ocean
  - prod DB
  - prod firebase project
  - prod email server
    - SendGrid or AWS SES
  - envs
  - ssl/tls
  - domena?
  - tags
  - zero downtime redeploy?
- predani
  - readme
    - functionality overview ✅ 
    - setup ✅
    - link to api docs ✅
    - architecture decisions? ✅
    - future improvements? ✅
    - popis flows


