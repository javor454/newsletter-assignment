FROM golang:1.23.1-alpine as base

RUN apk update \
  && apk add --no-cache bash~=5 \
  && apk add --no-cache make~=4 \
  && apk add --no-cache build-base~=0.5 \
  && apk add --no-cache gettext~=0.22 \
  && apk add --no-cache postgresql14-client~=14 \
  && apk add --no-cache --update-cache --upgrade curl~=8

ARG PROJECT_ROOT
ENV CGO_ENABLED=1 \
  GOROOT='/usr/local/go' \
  GO111MODULE='on' \
  PROJECT_ROOT=${PROJECT_ROOT} \
  DEBUG_DLV="0" \
  CONFIG_APP_ENV='prod'

WORKDIR ${PROJECT_ROOT}

COPY go.sum .
COPY go.mod .

RUN go mod download
RUN go mod tidy

COPY . .

##### development ##############################################################
FROM base as development

RUN curl -sSfL "https://raw.githubusercontent.com/cosmtrek/air/master/install.sh" | sh -s -- -b "$(go env GOPATH)/bin"

CMD ["/bin/sh", "-c", "air"]
#CMD ["sh", "-c", "${PROJECT_ROOT}/../build/newsletter-assignment"]
##### build ###########+#########################################################
FROM base as build
RUN bash script/build.sh

##### production ###############################################################
FROM alpine:3.20.3 as production

RUN apk update \
  && apk add --no-cache bash~=5 \
  && apk add --no-cache postgresql14-client~=14 \
  && apk add --no-cache curl~=8

ARG PROJECT_ROOT

RUN adduser -D -s /bin/sh -u 241 app

COPY --from=build "${PROJECT_ROOT}/build/newsletter-assignment" "/bin/newsletter-assignment"
USER app

CMD ["/bin/sh", "-c", "/bin/newsletter-assignment"]
