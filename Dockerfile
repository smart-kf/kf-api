#FROM golang:alpine as builder
#COPY . /app
#WORKDIR /app
#ENV GOPROXY=https://goproxy.io
#COPY . /app
#WORKDIR /app
#RUN go build -o app cmd/main.go
#
#FROM alpine
#WORKDIR /
#COPY --from=builder /app/app /

FROM alpine
RUN apk add --no-cache tzdata curl
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY data/static /app/static
COPY data/website /app/data/website
RUN curl -L -o /app/ip2region.xdb https://github.com/lionsoul2014/ip2region/raw/refs/heads/master/data/ip2region.xdb
COPY bin/app /app/