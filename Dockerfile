# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.20-alpine as builder

# 添加gcc环境 以支持sqlite
RUN apk --no-cache add gcc musl-dev

# 设置工作目录
WORKDIR /go/src/github.com/sxz799/fileshare-server

# 将应用的代码复制到容器中
COPY ./server/ .


# 编译应用程序
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env \
    && go mod tidy \
    && go build -o app .


FROM node:16


WORKDIR /gsCheck-web/
COPY ./web/ .


RUN npm config set registry https://registry.npm.taobao.org/ \
    && npm install \
    && npm run build


# 使用alpine镜像
FROM alpine:latest
# 添加bash
RUN apk add bash
# 配置时区
RUN apk update && apk add tzdata 
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime 
RUN echo "Asia/Shanghai" > /etc/timezone

WORKDIR /home

COPY --from=0 /go/src/github.com/sxz799/fileshare-server/app ./
COPY --from=0 /go/src/github.com/sxz799/fileshare-server/conf.yaml ./
COPY --from=1 /gsCheck-web/dist/ ./dist

RUN mkdir "cert"
RUN mkdir "conf"
RUN mkdir "files"

# 开放端口
EXPOSE 4000

# 运行应用程序
CMD ["./app"]
