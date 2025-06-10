FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /sharing-service

FROM gcr.io/distroless/base

COPY --from=builder /sharing-service /sharing-service

EXPOSE 3333

ENTRYPOINT ["/sharing-service"]