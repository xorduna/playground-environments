FROM --platform=linux/amd64 golang:1.21-alpine AS builder
ENV CGO_ENABLED 1
ARG BUILD_REF
ARG GOARCH

# Copy the source code into the container.
COPY . /attendance

# Build the service binary.
WORKDIR /attendance/apps/api

RUN apk add build-base
RUN GOARCH=${GOARCH} CGO_ENABLED=1 go build -o server -ldflags "-X main.build=${BUILD_REF}"

FROM --platform=linux/amd64 alpine:latest
#RUN apk --no-cache add tzdata

ARG BUILD_DATE
ARG BUILD_REF
ENV BUILD_REF=${BUILD_REF}

COPY --from=builder /attendance/apps/api/server /attendance/server
WORKDIR /attendance

EXPOSE 80

CMD ["./server"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="api" \
      org.opencontainers.image.vendor="Maulabs SL"
