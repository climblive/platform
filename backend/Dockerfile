FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./climblive ./cmd/api/main.go

ENTRYPOINT [ "./climblive" ]