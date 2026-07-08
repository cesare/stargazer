FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o stargazer cmd/server/main.go


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/stargazer .
COPY --chown=nobody:nobody --chmod=400 stargazer-config.toml /app/

USER nobody
CMD [ "./stargazer" ]
