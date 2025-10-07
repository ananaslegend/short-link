ARG CMD_NAME_ARG=api

FROM golang:1.24.3-alpine AS builder

ARG CMD_NAME_ARG
ENV CMD_NAME=${CMD_NAME_ARG}

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

WORKDIR /app/cmd/${CMD_NAME}
RUN go build -o /bin/${CMD_NAME}

FROM alpine:latest

ARG CMD_NAME_ARG
ENV CMD_NAME=${CMD_NAME_ARG}

COPY --from=builder /bin/${CMD_NAME} /app
COPY .env .

RUN chmod +x /app

CMD ["/app"]
