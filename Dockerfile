FROM golang:alpine

ENV NODE_ENV production

RUN apk -U add alpine-sdk libgit2-dev git gcc nodejs
RUN npm install -g yarn
RUN mkdir -p /go/src/github.com/checkr/codeflow/dashboard
COPY dashboard/yarn.lock /go/src/github.com/checkr/codeflow/dashboard/yarn.lock
COPY server/configs/codeflow.yml /etc/codeflow.yml
WORKDIR /go/src/github.com/checkr/codeflow/dashboard
RUN yarn
WORKDIR /go/src/github.com/checkr/codeflow/server
COPY . /go/src/github.com/checkr/codeflow
RUN go build -o /go/bin/codeflow .

EXPOSE 3000 3001 3002 9000
