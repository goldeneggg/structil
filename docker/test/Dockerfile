FROM golang:1.18-alpine

ENV CGO_ENABLED=1
ENV APP_DIR /go/src/github.com/goldeneggg/structil

RUN apk add make gcc g++ bash

WORKDIR ${APP_DIR}

COPY go.mod ${APP_DIR}
COPY go.sum ${APP_DIR}
RUN go get ./...
# COPY --from=structil/mod:latest /go/pkg/mod /go/pkg/mod

COPY . ${APP_DIR}

ENTRYPOINT ["make"]
CMD ["test"]
