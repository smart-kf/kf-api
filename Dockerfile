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
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY data/static /app/static
COPY data/website /app/data/website
COPY bin/app /app/