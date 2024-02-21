FROM golang:1-alpine AS build

RUN apk update && apk add make git gcc musl-dev

ADD . /app/src/wedding

WORKDIR /app/src/wedding

RUN make clean install
RUN make wedding

RUN mv ${SERVICE} /${SERVICE}

FROM alpine:latest

ARG SERVICE

RUN apk add --no-cache ca-certificates && mkdir /app
COPY --from=build /${SERVICE} /app/${SERVICE}

ENTRYPOINT ["/app/energy-industry-actuals-repository"]
