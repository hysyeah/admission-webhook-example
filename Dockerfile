FROM golang:1.18-alpine

WORKDIR /app

COPY . .
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o main

RUN chmod +x main

EXPOSE 443
CMD ["./main"]