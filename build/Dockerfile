FROM golang:alpine as builder
RUN apk --update --no-cache add bash
WORKDIR /app
ADD . .
RUN go build -o app ./cmd

FROM alpine as prod
WORKDIR /app
COPY --from=builder /app/app /app/app
COPY --from=builder /app/config.yml /app/config.yml
COPY --from=builder /app/index.txt /app/index.txt
EXPOSE 8080
CMD ["./app"]
