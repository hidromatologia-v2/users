FROM golang:1.21-alpine AS build-stage

RUN apk add upx
WORKDIR /users-src
COPY . .
RUN go build -o /users .
RUN upx /users

FROM alpine:latest AS release-stage

COPY --from=build-stage /users /users
# -- Environment variables
ENV MEMPHIS_STATION     "messages"
ENV MEMPHIS_PRODUCER    "messages-producer"
ENV MEMPHIS_HOST        "memphis"
ENV MEMPHIS_USERNAME    "root"
ENV MEMPHIS_PASSWORD    "memphis"
ENV MEMPHIS_CONN_TOKEN  ""
ENV JWT_SECRET          "SECRET"
ENV REDIS_EMAIL_ADDR    "redis:6379"
ENV REDIS_EMAIL_DB      "1"
ENV REDIS_PASSWORD_ADDR "redis:6379"
ENV REDIS_PASSWORD_DB   "2"
ENV POSTGRES_DSN        "host=postgres user=sulcud password=sulcud dbname=sulcud port=5432 sslmode=disable"
# -- Environment variables
ENTRYPOINT [ "sh", "-c", "/users :5000" ]