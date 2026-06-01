FROM golang:1.26.2-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
COPY db/ db/
COPY handlers/ handlers/
COPY models/ models/
RUN go mod tidy
RUN CGO_ENABLED=1 go build -o server .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/server /usr/local/bin/server
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8080
CMD ["server"]
