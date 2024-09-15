# Newsletter assignment

## Development plan
- Setup project via docker compose
  - go app
  - postgre
- Setup other dependencies ??????
  - firebase 
    - setup project and configure SDK
  - email server
    - SendGrid or AWS SES
- MVP go app
  - REST
    - user registration
      - db storing password
    - user authentication JWT
    - user authorization
    - newsletter management
    - subscription management
    - email functionality
    - api documentation - swagger
    - healthcheck
    - graceful shutdown WIP
    - logs
    - panic recovery
    - security
      - rate limiting ????
      - cors
  - testing
    - system tests (func, unit, integration)
- deployment
  - digital ocean
  - prod DB
  - prod firebase project
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
- EXTRA
  - basic performance test 
  - CI (github actions lint / build / deploy)
  - simple backoffice
  - graphql

## Architecture
- auth component
  - functions
    - registration
    - authorization
    - authentication
  - requirements
    - 100 users
- news component
  - functions
    - newsletter management
    - subscription management
    - email management
  - requirements

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
