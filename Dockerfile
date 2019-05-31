FROM golang:1.11-alpine as build

# Install tools required for project
# Run `docker build --no-cache .` to update dependencies
RUN apk add --no-cache build-base git
ENV APP_PATH /go/src/github.com/raazcrzy/imdb/app
WORKDIR  ${APP_PATH}

# Add eaas-provisioner files
ADD .  /go/src/github.com/raazcrzy/imdb

RUN go install
# Define default command
CMD ["/go/bin/app"]

# Expose Ports
EXPOSE 8080