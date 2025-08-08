FROM golang:1.24.4-alpine as builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o avtorizate ./cmd

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/avtorizate .

EXPOSE 50053

CMD ["./avtorizate"]