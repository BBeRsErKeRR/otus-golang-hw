
ARG bin_file /opt/calendar/calendar-app
ARG src_main cmd/calendar/*
ARG src_config ./configs/config.toml

# Собираем в гошке
FROM golang:1.19 as build

ENV BIN_FILE=${bin_file} \
    CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} ${src_main}

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE=${bin_file}
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/config.toml
COPY ${src_config} ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
