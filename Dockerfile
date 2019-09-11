# syntax=docker/dockerfile:1.0.0-experimental 
FROM golang:1.12.9-alpine as build

RUN apk add --no-cache openssh-client git ca-certificates

WORKDIR /usr/src/git-app

RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

RUN git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"

COPY go.mod go.sum ./

RUN --mount=type=ssh go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o go-git-api

FROM scratch

COPY --from=build /usr/src/git-app/go-git-api go-git-api

CMD [ "./go-git-api" ]
