FROM archlinux/base

ARG NODE_ENV=production
ENV NODE_ENV=${NODE_ENV}
ENV APP_PATH /go/src/github.com/checkr/codeflow
ENV GOPATH /go
ENV PATH ${PATH}:/go/bin

COPY docker-entrypoint.sh $APP_PATH/docker-entrypoint.sh

# Go application cache buster
COPY server $APP_PATH/server

# Base updates
RUN pacman -Sy archlinux-keyring --noconfirm && \
	pacman -Syu --noconfirm && \
	rm -rf /var/lib/pacman/pkg/*

RUN pacman -Sy --noconfirm libgit2 git openssh gcc go go-tools base-devel

# create .ssh direcotry so git can create known_hosts file
RUN mkdir ~/.ssh

# development dependencies
RUN go get github.com/cespare/reflex

# Go application
WORKDIR $APP_PATH/server
RUN go build -i -o /go/bin/codeflow .

# Dashboard cache buster
COPY dashboard $APP_PATH/dashboard

# Node dependencies
RUN pacman -Sy --noconfirm nodejs npm
WORKDIR $APP_PATH/dashboard
RUN npm install
RUN npm run build

# Docs
WORKDIR $APP_PATH/docs
COPY docs $APP_PATH/docs
RUN npm install
RUN npm install gitbook-cli -g
RUN gitbook install && gitbook build

WORKDIR $APP_PATH

ENTRYPOINT ["./docker-entrypoint.sh"]
