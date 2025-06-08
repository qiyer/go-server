FROM golang:1.24

# 设置 GOPROXY 环境变量
ENV GOPROXY=https://goproxy.cn

RUN mkdir /app

ADD . /app

WORKDIR /app
# 复制所有必要文件（确保包含 data/config.json）
COPY . .   

RUN go build -o main cmd/main.go

CMD ["/app/main"]