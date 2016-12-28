FROM golang:alpine

ENV NODE_ENV production

RUN apk -U add alpine-sdk libgit2-dev git gcc nodejs
RUN npm install -g yarn
COPY . /go/src/github.com/checkr/codeflow
COPY server/configs/codeflow.yml /etc/codeflow.yml
RUN cd /go/src/github.com/checkr/codeflow/server && go build -o /go/bin/codeflow .

WORKDIR /go/src/github.com/checkr/codeflow/client
RUN yarn
RUN npm run build

EXPOSE 3000 3001 3002 9000
