FROM golang:1.16.13-alpine3.15 AS build

ARG APP_NAME=httpserver
WORKDIR /go/src/project
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o ${APP_NAME}

FROM busybox
COPY --from=build /go/src/project/${APP_NAME} /app/
ENTRYPOINT ["/app/httpserver"]

# 构建命令
# cd cn-camp/homework/m02
# docker build -f ./.docker/Dockerfile -t httpserver:v0.1 .
