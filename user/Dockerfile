FROM docker.io/library/golang:1-alpine AS build
RUN apk update && apk upgrade && apk add --no-cache gcc pkgconfig make jq tar xz czmq libzmq musl-dev linux-headers
WORKDIR /app
COPY . .
ENV CGO_ENABLED=1
RUN make build
RUN set -ex && \
    wget https://github.com/upx/upx/releases/download/v4.0.2/upx-4.0.2-amd64_linux.tar.xz && \
    tar -xvf upx-4.0.2-amd64_linux.tar.xz upx-4.0.2-amd64_linux/upx && \
    mv ./upx-4.0.2-amd64_linux/upx . && \
    ./upx --no-color --mono --no-progress --ultra-brute --no-backup ./bin/user && \
    ./upx --test ./bin/user
FROM gcr.io/distroless/base-debian11:nonroot
COPY --from=build --chown=nonroot:nonroot /app/bin/user user
ENV TZ=UTC
ENTRYPOINT [ "./user" ]
