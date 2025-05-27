FROM golang:1.24

RUN mkdir /app

ADD . /app

WORKDIR /app
# 复制所有必要文件（确保包含 data/config.json）
COPY . .   

RUN go build -o main cmd/main.go

CMD ["/app/main"]