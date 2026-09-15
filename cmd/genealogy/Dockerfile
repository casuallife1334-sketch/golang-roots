FROM golang:1.25.5-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/genealogy ./cmd/genealogy

FROM alpine:3.24.0

WORKDIR /app

COPY --from=builder /app/genealogy /app/genealogy

ENV PROJECT_ROOT=/app

CMD ["/app/genealogy"]
