
# 基础镜像
FROM golang:1.21.1

ADD . /shorturl

# 设置工作目录
WORKDIR /shorturl


ENV GOPROXY=https://goproxy.cn,direct


# 构建项目
RUN go build -o shorturl-url ./rpc/url 
RUN go build -o shorturl-user ./rpc/user 
RUN go build -o shorturl-client .

# 暴露服务端口
EXPOSE 8888
EXPOSE 8889
EXPOSE 8080

CMD  ./shorturl-user & ./shorturl-url & ./shorturl-client