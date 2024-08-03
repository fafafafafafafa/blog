# FROM golang:latest
FROM golang:1.20 AS builder

# ENV GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=https://mirrors.aliyun.com/goproxy,direct
ENV GO111MODULE=on
# 用来控制golang 编译期间是否支持调用 cgo 命令的开关
# 当CGO_ENABLED=1， 进行编译时， 会将文件中引用libc的库（比如常用的net包）以动态链接的方式生成目标文件
# 当CGO_ENABLED=0， 进行编译时， 则会把在目标文件中未定义的符号（外部函数）一起链接到可执行文件中
ENV CGO_ENABLED=0 


WORKDIR /go-gin-example
# COPY . /go-gin-example
COPY . .

# RUN curl -ik https://goproxy.cn/golang.org/x/text/@latest
# RUN go env -w GOPROXY=https://goproxy.cn
RUN go mod download

RUN go build -o main .

# ENTRYPOINT ["./go-gin-example"]

###################
# 接下来创建一个小镜像
###################
FROM scratch

WORKDIR /go-gin-example

# 拷贝配置文件文件
COPY ./conf ./conf

EXPOSE 8000
# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /go-gin-example/main .

ENTRYPOINT ["/go-gin-example/main"]
