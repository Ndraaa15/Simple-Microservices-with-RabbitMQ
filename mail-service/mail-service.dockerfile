FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN ls -al

RUN CGO_ENABLED=0 go build -o mailApp ./cmd/api

FROM alpine:latest 

RUN apk --no-cache add curl

WORKDIR /app

COPY --from=builder /app/mailApp .
COPY --from=builder /app/templates /app/templates

ENTRYPOINT ["./mailApp"]
