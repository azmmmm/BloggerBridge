FROM golang:1.20 AS build-stage
WORKDIR /code
COPY go.mod go.sum ./

RUN go env -w GOPROXY=https://goproxy.cn
RUN\
    --mount=type=cache,target=/var/cache/golang \
    go mod download



COPY ./src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -o ./wrapper ./src/main.go



FROM gcr.io/distroless/base-debian11:latest AS build-release-stage

WORKDIR /code

COPY --from=build-stage /code/wrapper ./wrapper
COPY ./config ./config

ENV TZ=Asia/Shanghai

ENTRYPOINT [ "./wrapper" ]

EXPOSE 8085