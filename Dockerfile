# Dockerfile
FROM golang:1.24.3-alpine

WORKDIR /app

COPY . .

WORKDIR /app/apps/organization/database
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /db_init

# Compile le fichier db_init.go
# RUN go build -o /db_init

CMD ["/db_init"]
