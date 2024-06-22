FROM golang:1.22.4 as BuildStage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /shortener cmd/shortener/main.go

FROM alpine:latest

WORKDIR /
COPY --from=BuildStage /shortener /shortener

CMD ["/shortener"]