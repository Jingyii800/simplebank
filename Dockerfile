FROM golang:1.21.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY db /app/db

EXPOSE 8080 8888
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]