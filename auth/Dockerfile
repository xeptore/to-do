FROM docker.io/library/golang:1-alpine AS build
RUN apk update && apk upgrade && apk add --no-cache make jq tar xz musl-dev linux-headers git
WORKDIR /app
ENV GOPRIVATE=github.com/xeptore/to-do
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN make build
RUN set -ex && \
    wget https://github.com/upx/upx/releases/download/v4.0.2/upx-4.0.2-amd64_linux.tar.xz && \
    tar -xvf upx-4.0.2-amd64_linux.tar.xz upx-4.0.2-amd64_linux/upx && \
    mv ./upx-4.0.2-amd64_linux/upx . && \
    ./upx --no-color --mono --no-progress --ultra-brute --no-backup ./bin/auth && \
    ./upx --test ./bin/auth
FROM scratch
COPY --from=build /app/bin/auth /auth
ENV TZ=UTC
ENTRYPOINT [ "/auth" ]
