# Собираем в гошке
FROM golang:1.19 as build

ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY migrations/. .

CMD ["goose", "up"]
