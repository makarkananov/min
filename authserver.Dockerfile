FROM golang:1.22.4 as BuildStage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth cmd/auth/main.go

FROM alpine:latest

WORKDIR /
COPY --from=BuildStage /auth /auth

CMD ["/auth"]