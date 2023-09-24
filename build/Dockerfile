FROM golang:1.21 as build

ENV BIN_FILE /opt/banner-rotation/banner-rotation
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -o ${BIN_FILE} cmd/*

FROM alpine:3.9

ENV BIN_FILE "/opt/banner-rotation/banner-rotation"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV LOG_LEVEL 4
ENV CONFIG_FILE /etc/banner-rotation/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}
COPY ./internal/storage/migrations/ /opt/banner-rotation

CMD ${BIN_FILE} -config ${CONFIG_FILE} -logLevel ${LOG_LEVEL}
