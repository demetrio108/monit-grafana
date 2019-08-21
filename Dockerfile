FROM golang:1.12-alpine as builder

MAINTAINER demetrio108 <demetrio108@protonmail.com>

ENV GODEBUG netdns=cgo
RUN apk add --no-cache --update alpine-sdk git make

COPY . /go/src/github.com/demetrio108/monit-grafana
WORKDIR /go/src/github.com/demetrio108/monit-grafana
RUN go get
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s -w -extldflags "-static"' .
RUN mkdir -p rootfs/etc && \
    cp monit-grafana rootfs/ && \
    cp monit-grafana.yml rootfs/etc/

FROM scratch

COPY --from=builder /go/src/github.com/demetrio108/monit-grafana/rootfs/ /

EXPOSE 8080

CMD ["/monit-grafana", "-c", "/etc/monit-grafana.yml"]
