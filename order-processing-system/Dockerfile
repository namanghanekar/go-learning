FROM alpine:3.21
RUN adduser -D app
USER app
ARG SERVICE
COPY bin/${SERVICE} /service
ENTRYPOINT ["/service"]
