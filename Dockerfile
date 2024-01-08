FROM golang:1.21.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db /app/db

# Make sure the script is executable
RUN chmod +x /app/start.sh

EXPOSE 8080 8888
ENTRYPOINT [ "/app/start.sh" ]
CMD [ "/app/main" ]