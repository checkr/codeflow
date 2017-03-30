FROM golang:alpine

ARG NODE_ENV=production
ENV NODE_ENV=${NODE_ENV}
ENV APP_PATH /go/src/github.com/checkr/codeflow

RUN mkdir -p $APP_PATH
WORKDIR $APP_PATH

RUN apk -U add alpine-sdk libgit2-dev git gcc nodejs
COPY ./dashboard/package.json $APP_PATH/dashboard/package.json
COPY ./server/configs/codeflow.yml /etc/codeflow.yml
RUN cd $APP_PATH/dashboard/ && npm install
COPY . /go/src/github.com/checkr/codeflow

WORKDIR $APP_PATH/server
RUN go build -i -o /go/bin/codeflow .
RUN go get github.com/cespare/reflex

WORKDIR $APP_PATH/dashboard
RUN npm run build

RUN npm install gitbook-cli -g
RUN cd $APP_PATH/docs && npm install && gitbook install && gitbook build

WORKDIR $APP_PATH
