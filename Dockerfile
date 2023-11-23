FROM golang:alpine as builder

COPY . $GOPATH/src/gin-dbo
WORKDIR $GOPATH/src/gin-dbo

COPY . .
RUN go mod tidy
RUN apk add --no-cache build-base

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o /gin-dbo .

FROM alpine:3.4

COPY --from=builder /gin-dbo /gin-dbo
COPY .env .env
RUN touch .env

EXPOSE 30001

ENTRYPOINT ["./gin-dbo"]