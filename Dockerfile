# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o network .

# --- Runtime image ---
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/network .
COPY migrations ./migrations
ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=yourpassword
ENV DB_NAME=social_network
ENV JWT_SECRET=your_jwt_secret
ENV PORT=8080

CMD ["./network"] 