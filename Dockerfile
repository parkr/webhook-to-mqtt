FROM golang:1.22.4-bullseye as builder
WORKDIR /workspace
EXPOSE 3306
COPY . .
RUN go version
RUN set -ex \
  && go install github.com/parkr/webhook-to-mqtt/... \
  && ls -l /go/bin/webhook-to-mqtt

FROM debian:bullseye-slim
COPY --from=builder /go/bin/* /go/bin/
ENTRYPOINT [ "/go/bin/webhook-to-mqtt" ]