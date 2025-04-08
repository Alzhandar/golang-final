FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

FROM alpine:latest  

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D -g '' appuser

ENV TZ=Asia/Almaty

WORKDIR /app

COPY --from=builder /app/main /app/
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/migrations /app/migrations

USER appuser

EXPOSE 8080

CMD ["./main"]