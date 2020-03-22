FROM golang:1.13

ENV CGO_ENABLED=1
ENV APP_DIR /go/src/github.com/goldeneggg/structil

WORKDIR ${APP_DIR}

COPY go.mod ${APP_DIR}
COPY go.sum ${APP_DIR}
RUN GO111MODULE=on go mod download

COPY . ${APP_DIR}

ENTRYPOINT ["make"]
CMD ["test"]
