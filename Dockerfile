FROM finalduty/archlinux:monthly

ARG NODE_ENV=production
ENV NODE_ENV=${NODE_ENV}
ENV APP_PATH /go/src/github.com/checkr/codeflow
ENV GOPATH /go
ENV PATH ${PATH}:/go/bin

RUN mkdir -p $APP_PATH
WORKDIR $APP_PATH

RUN pacman -Syu --noconfirm
RUN pacman -Sy --noconfirm libgit2 git gcc nodejs go go-tools npm base-devel

RUN go get github.com/cespare/reflex
RUN npm install gitbook-cli -g

WORKDIR $APP_PATH/dashboard
COPY ./dashboard/package.json ./package.json
RUN npm install

WORKDIR $APP_PATH/docs
COPY ./docs/package.json ./package.json
RUN npm install

COPY . $APP_PATH

WORKDIR $APP_PATH/server
RUN go build -i -o /go/bin/codeflow .

WORKDIR $APP_PATH/docs
RUN gitbook install && gitbook build

WORKDIR $APP_PATH/dashboard
RUN npm run build

WORKDIR $APP_PATH

ENTRYPOINT ["./docker-entrypoint.sh"]
