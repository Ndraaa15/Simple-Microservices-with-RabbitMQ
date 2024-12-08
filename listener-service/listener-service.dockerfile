FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o listenerApp .

FROM alpine:latest 

RUN apk --no-cache add curl

WORKDIR /app

COPY --from=builder /app/listenerApp .

ENTRYPOINT ["./listenerApp"]
