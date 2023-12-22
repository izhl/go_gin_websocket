#   编译阶段
FROM golang

#   声明参数
ARG GO_ENVIRONMENT

#   定义环境变量
ENV GO111MODULE=on \
    GOPROXY="https://goproxy.cn,direct" \
    GOSUMDB="sum.golang.google.cn" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

#   设置工作目录
WORKDIR /go/src/go_gin_websocket/

#   拷贝项目文件
COPY . .

#   开始编译
RUN  go build

#   运行阶段
FROM alpine:3

#   设置工作目录
WORKDIR /app

#   拷贝编译后可运行文件
COPY --from=0 /go/src/go_gin_websocket .

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk --no-cache add tzdata \
        bash \
    && chmod 0777 ./is_run.sh

#   运行服务
CMD ["./go_gin_websocket"]