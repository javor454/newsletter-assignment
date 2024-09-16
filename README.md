# Newsletter assignment

## Development plan
- Setup project via docker compose ✅
  - go app ✅
  - postgre ✅
  - firebase ✅
    - configure SDK
  - email server
- air hot rebuild in local env ✅
- Setup other dependencies
  - email server
    - SendGrid or AWS SES
- db migrations
  - structures  ✅
  - data 
- MVP go app
  - REST
    - user registration ✅
      - db storing password ✅ 
    - user authentication JWT ✅
    - user authorization ✅
    - newsletter management
    - subscription management
    - email functionality
    - api documentation - swagger
    - healthcheck
    - graceful shutdown
    - logs
    - panic recovery
    - security
      - rate limiting ????
      - cors
    - pagination
  - testing
    - system tests (func, unit, integration)
- FIX
  - zkontrolovat volumes
- deployment
  - digital ocean
  - prod DB
  - prod firebase project
  - prod email server
    - SendGrid or AWS SES
  - envs
  - ssl/tls
  - domena?
- predani
  - readme
    - functionality overview
    - setup 
    - link to api docs
    - architecture decisions?
    - future improvements?
    - popis flows
- EXTRA
  - basic performance test
  - CI (github actions lint / build / deploy)
  - simple backoffice
  - graphql
- LINKY
  Editor creates a newsletter
  System generates a unique subscription link for the newsletter
  Editor shares the unique subscription link
  User clicks on the unique link
  System verifies the link's validity
  User is presented with a form to enter their email
  User submits their email
  System checks if the email is already subscribed to this newsletter
  If not subscribed, system adds the email to Firebase as a new subscriber
  System sends a confirmation email to the subscriber

## Data layer
- pagination
- postgres
- user
CREATE TABLE users (
id SERIAL PRIMARY KEY,
email VARCHAR(255) UNIQUE NOT NULL,
password_hash VARCHAR(255) NOT NULL,
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
- news
CREATE TABLE newsletters (
id SERIAL PRIMARY KEY,
firebase_id VARCHAR(255) UNIQUE NOT NULL,
user_id INTEGER REFERENCES users(id),
name VARCHAR(255) NOT NULL,
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
- links
  CREATE TABLE subscription_links (
  id SERIAL PRIMARY KEY,
  newsletter_id INTEGER REFERENCES newsletters(id),
  unique_token VARCHAR(64) UNIQUE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP WITH TIME ZONE
  );


- firebase = realtime database
subscriptions/
├── {subscriberEmail}/
│   ├── {newsletterId1}: true
│   ├── {newsletterId2}: true
│   └── ...

newsletters/
├── {newsletterId}/
│   ├── name: string
│   ├── description: string
│   ├── createdAt: timestamp
