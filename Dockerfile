FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/wn ./cmd/app

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/wn .
COPY config/main.yaml config/main.yaml
COPY migrations/ migrations/
COPY docs/ docs/
COPY statics/ statics/

EXPOSE 8080

CMD ["./wn"]
