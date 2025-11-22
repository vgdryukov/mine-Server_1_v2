FROM golang:1.21-alpine as builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . ./
RUN go build -o cafe-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/cafe-server .

EXPOSE 8080
CMD ["./cafe-server"]