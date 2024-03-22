## Build
FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
ADD constants /app/constants
ADD github /app/github

RUN go build -o /raito-cli

## Deploy
FROM alpine:3 as deploy

LABEL org.opencontainers.image.base.name="alpine:3"

RUN apk add --no-cache tzdata

WORKDIR /app

RUN mkdir -p /config

ENV TZ=Etc/UTC
ENV CLI_CRON="0 2 * * *"

RUN addgroup -S raito && adduser -D -S -G raito --no-create-home raito && chmod +w /tmp
RUN chown raito:raito /app /config

COPY --from=build /raito-cli /raito
RUN chown raito:raito /raito

USER raito

ENTRYPOINT /raito run -c "$CLI_CRON" --config-file /config/raito.yml --log-output

## Deploy-amazon
FROM amazon/aws-cli:2.15.31 as amazonlinux

LABEL org.opencontainers.image.base.name="amazon/aws-cli:2.15.10"

RUN yum -y install tzdata jq shadow-utils

WORKDIR /app

RUN mkdir -p /config

ENV TZ=Etc/UTC
ENV CLI_CRON="0 2 * * *"

RUN groupadd -r raito && useradd -r -g raito raito
RUN chown raito:raito /app /config

COPY --from=build /raito-cli /raito

RUN chown raito:raito /raito
USER raito

ENTRYPOINT []
CMD /raito run -c "$CLI_CRON" --config-file /config/raito.yml --log-output
