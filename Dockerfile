FROM golang:1.21-alpine as build
RUN apk add -u git
WORKDIR /app
COPY pkg/ pkg/
COPY cmd/ cmd/
COPY main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN go build -o /caretaker

FROM alpine
RUN apk add -u ca-certificates
COPY --from=build /caretaker /app/

LABEL "name"="Caretaker"
LABEL "maintainer"="Gergely Brautigam <gergely@gergelybrautigam.com>"
LABEL "version"="0.0.1"

LABEL "com.github.actions.name"="Caretake - Project Manager"
LABEL "com.github.actions.description"="Manage project issues automatically."
LABEL "com.github.actions.icon"="package"
LABEL "com.github.actions.color"="purple"

WORKDIR /app/
ENTRYPOINT [ "/app/caretaker" ]
