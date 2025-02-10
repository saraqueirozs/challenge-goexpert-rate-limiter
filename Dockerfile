FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./cmd

FROM alpine:latest
COPY --from=builder /app .
CMD [ "./server" ]
