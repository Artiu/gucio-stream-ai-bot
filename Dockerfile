FROM golang:1.22.0 AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /gucio-stream-ai-bot

FROM alpine
WORKDIR /
COPY --from=base /gucio-stream-ai-bot /gucio-stream-ai-bot
RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Warsaw /etc/localtime
CMD ["/gucio-stream-ai-bot"]
