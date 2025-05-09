# syntax=docker/dockerfile:1

FROM golang:1.24-alpine

WORKDIR /app

COPY tmate-server /app/service
COPY assets/ /app/assets
RUN chmod +x /app/service
RUN mkdir /app/logs

ENV TMATE_PORT=8080
ENV TZ=Europe/Berlin

RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Berlin /etc/localtime

EXPOSE 8080
EXPOSE 465
EXPOSE 587

ENTRYPOINT [ "./service" ]
