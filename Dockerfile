FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /tg ./cmd/main/main.go

FROM alpine:latest
COPY --from=builder /tg /tg
RUN chmod +x /tg
EXPOSE 8080
CMD ["/tg"]
