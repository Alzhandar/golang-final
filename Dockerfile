FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/local/bin/migrate

FROM alpine:latest  

RUN apk --no-cache add ca-certificates tzdata netcat-openbsd

WORKDIR /app

COPY --from=builder /app/main /app/
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /usr/local/bin/migrate /usr/bin/migrate

COPY entrypoint.sh /app/
RUN chmod +x /app/entrypoint.sh

RUN adduser -D -g '' appuser
ENV TZ=Asia/Almaty

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./main"]