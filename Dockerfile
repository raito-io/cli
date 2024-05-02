## Build
FROM golang:1.22-alpine AS build
ARG VERSION
ARG COMMIT_DATE

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY Makefile ./
COPY *.yaml ./
COPY *.yml ./
ADD base /app/base
ADD cmd /app/cmd
ADD internal /app/internal
ADD proto /app/proto
ADD scripts /app/scripts

RUN apk add --no-cache make
RUN go install github.com/bufbuild/buf/cmd/buf@v1.30.0
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

RUN make generate
RUN go build -o raito main.go -ldflags="-X main.version=$VERSION -X main.date=$COMMIT_DATE"

## Deploy
FROM alpine:3 as deploy

LABEL org.opencontainers.image.base.name="alpine:3"

RUN apk add --no-cache tzdata

RUN mkdir -p /config

ENV TZ=Etc/UTC
ENV CLI_CRON="0 2 * * *"

RUN addgroup -S raito && adduser -D -S -G raito raito && chmod +w /tmp
RUN chown raito:raito /config

COPY --from=build /app/raito /raito
RUN chown raito:raito /raito

USER raito

ENTRYPOINT /raito run -c "$CLI_CRON" --config-file /config/raito.yml --log-output

## Deploy-amazon
FROM amazon/aws-cli:2.15.31 as amazonlinux

LABEL org.opencontainers.image.base.name="amazon/aws-cli:2.15.10"

RUN yum -y install tzdata jq shadow-utils

RUN mkdir -p /config

ENV TZ=Etc/UTC
ENV CLI_CRON="0 2 * * *"

RUN groupadd -r raito && useradd -r -g raito raito
RUN chown raito:raito /config

COPY --from=build /app/raito /raito

RUN chown raito:raito /raito
USER raito

ENTRYPOINT []
CMD /raito run -c "$CLI_CRON" --config-file /config/raito.yml --log-output
