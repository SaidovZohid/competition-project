FROM golang:1.19.1-alpine3.16 as builder

WORKDIR /url-shorter

COPY . .

RUN go build -o main cmd/main.go

FROM alpine:3.16

WORKDIR /url-shorter
RUN mkdir media

COPY --from=builder /url-shorter/main .
COPY templates ./templates

EXPOSE 8080

CMD [ "/url-shorter/main" ]