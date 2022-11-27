

# build stage
FROM golang:1.19-alpine AS build

ADD . /go/build
WORKDIR /go/build

# install gcc
RUN apk add build-base

# final stage
FROM alpine:3.16

RUN go build -o twss main.go

# install timezone data
RUN apk add --no-cache tzdata
ENV TZ Asia/Taipei
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=build /go/build/twss /var/application/twss
COPY --from=build /go/build/config /var/application/config

WORKDIR /var/application
CMD [ "./twss" ]