FROM golang:1.19 as build

ARG MAIN_SRC
ARG BIN_FILE
ENV BIN_FILE=${BIN_FILE} \
    MAIN_SRC=${MAIN_SRC} \
    CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
ARG BIN_FILE
ARG MAIN_SRC
RUN CGO_ENABLED=0 go build \
    -ldflags "$LDFLAGS" \
    -o ${BIN_FILE} ${MAIN_SRC}

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ARG BIN_FILE
ENV BIN_FILE=${BIN_FILE}
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ARG CONFIG_SRC
ENV CONFIG_FILE /etc/calendar/config.toml
COPY ${CONFIG_SRC} ${CONFIG_FILE}

CMD ${BIN_FILE} --config ${CONFIG_FILE}
