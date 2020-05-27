FROM golang:1.14-alpine AS rushb

WORKDIR /app

COPY . /app

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -mod vendor -gcflags "-l -w" -o appsrv

FROM alpine:3.11

LABEL maintainer="unknow"
LABEL k8s-app="appsrv"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update && apk --no-cache add tzdata ca-certificates wget \
    && cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

USER 1000
WORKDIR /app
COPY --from=rushb /app/resources/fonts/Songti.ttc /app/Songti.ttc
COPY --from=rushb /app/appsrv /app/appsrv

EXPOSE 8080

CMD [ "/app/appsrv" ]

ENTRYPOINT [ "/app/appsrv" ]
