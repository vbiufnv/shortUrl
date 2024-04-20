
# 基础镜像
FROM golang:1.21.1

ADD . /shorturl

WORKDIR /shorturl

# COPY go.mod .
# COPY go.sum .
# RUN go mod download


COPY . .



ENV GOPROXY=https://goproxy.cn,direct


# 构建项目
RUN go build -o shorturl-url ./rpc/url 
RUN go build -o shorturl-user ./rpc/user 
RUN go build -o shorturl-client .

# 暴露服务端口
EXPOSE 8888
EXPOSE 8889
EXPOSE 8080

# 只能有一个CMD指令 多个会覆盖

# ENTRYPOINT [ "" ]  在run之前执行  可以执行一个脚本 脚本中启动url user服务


CMD  ./shorturl-user & ./shorturl-url & ./shorturl-client

# RUN VS CMD      前者在镜像构建过程中执行  后者构建时不执行，只是配置镜像的默认程序入口