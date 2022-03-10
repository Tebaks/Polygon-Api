FROM golang:1.17 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /polygonapp

FROM alpine:3.11

WORKDIR /

COPY --from=builder /polygonapp /polygonapp
COPY --from=builder /app/config/config.yaml ./config/config.yaml

ENTRYPOINT [ "/polygonapp" ]