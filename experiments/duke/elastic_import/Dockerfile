FROM golang:1.11.4-alpine

RUN apk update && apk add --no-cache git gcc build-base

RUN mkdir /app

ADD . /app/

WORKDIR /app/experiments/duke/elastic_import

RUN ./build.sh

RUN adduser -S -D -H -h /app gouser

RUN chmod go+x ./cmd/widgets_import/widgets_import

RUN chmod go+x ./cmd/staging_import/staging_import

RUN chmod go+x ./cmd/elastic_import/elastic_import

WORKDIR /app/

RUN chmod go+x ./cmd/graphql_server/graphql_server

EXPOSE 9001

CMD ["./cmd/graphql_server/graphql_server"]

USER gouser
