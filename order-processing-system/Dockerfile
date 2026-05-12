# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
ARG SERVICE
COPY internal internal
COPY gen gen
COPY pkg pkg
COPY cmd/${SERVICE} cmd/${SERVICE}
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/${SERVICE}

FROM alpine:3.21
RUN adduser -D app
USER app
COPY --from=build /out/service /service
ENTRYPOINT ["/service"]
