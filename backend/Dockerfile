FROM golang:1.26

WORKDIR /app

ARG VERSION=dev

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/climblive/platform/backend/internal/handlers/rest.Version=${VERSION}" -o ./climblive ./cmd/api/main.go

ENTRYPOINT [ "./climblive" ]
