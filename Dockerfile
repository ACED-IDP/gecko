FROM golang:1.22.6-alpine AS build-deps
RUN apk add make git bash build-base libc-dev binutils-gold curl postgresql-client	

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOPATH=/go
ENV PATH="/go/bin:${PATH}"

WORKDIR $GOPATH/src/github.com/ACED-IDP/chameleon/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GITCOMMIT=$(git rev-parse HEAD) \
    GITVERSION=$(git describe --always --tags) \
    && go build \
    -ldflags="-X 'github.com/ACED-IDP/chameleon/chameleon/version.GitCommit=${GITCOMMIT}' -X 'github.com/ACED-IDP/chameleon/chameleon/version.GitVersion=${GITVERSION}'" \
    -o bin/chameleon

CMD ["bin/chameleon"]