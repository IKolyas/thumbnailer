FROM golang:1.23.3-alpine AS build

ENV BIN_FILE=/opt/previewer
ENV CODE_DIR=/go/src

RUN apk add --update --no-cache --repository git build-base glib-dev vips-dev

WORKDIR ${CODE_DIR}

COPY ./go.mod ${CODE_DIR}/go.mod
COPY ./go.sum ${CODE_DIR}/go.sum
RUN go mod download

COPY . ${CODE_DIR}

RUN go build -o ${BIN_FILE} cmd/previewer/main.go

FROM alpine:latest

RUN apk add --no-cache vips-dev

WORKDIR /app

RUN mkdir -p /app/logs && \
    touch /app/logs/previewer.log && \
    chmod 666 /app/logs/previewer.log

ENV BIN_FILE=/app/bin/previewer
COPY --from=build /opt/previewer ${BIN_FILE}

ENV CONFIG_FILE=/app/configs/config.json
COPY ./configs/config.json ${CONFIG_FILE}

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["${BIN_FILE} --config=${CONFIG_FILE}"]