# Dockerfile 
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN go build -o userservice ./cmd/userservice

# Run stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/userservice .

EXPOSE 8080

CMD [ "./userservice" ]