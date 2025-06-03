# Build stage
# 替换为zrjs这边的hub库
FROM registry.cn-shenzhen.aliyuncs.com/blin/go-base:0.1 AS builder
COPY . .
RUN go build -o main

# Run stage
#FROM golang:1.21-alpine
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/swagger.yaml .
COPY --from=builder /app/config.json .
EXPOSE 8081
EXPOSE 9001
CMD [ "/app/main" "-c" /app/config.json]
