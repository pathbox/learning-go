FROM kong/go-plugin-tool:latest-centos-7 AS builder

ENV GOPROXY https://goproxy.cn,direct

RUN mkdir -p /tmp/go/src/rate-limiting/

COPY . /tmp/go/src/rate-limiting

RUN cd /tmp/go/src/rate-limiting && \
    go build -buildmode plugin custom-rate-limiting.go && \
    cd /tmp/go/src/rate-limiting/go-pluginserver && \
    go build github.com/Kong/go-pluginserver

FROM kong:2.1.3-centos

ENV KONG_DATABASE off
ENV KONG_GO_PLUGINS_DIR /etc/kong/plugins/
#ENV KONG_DECLARATIVE_CONFIG /etc/kong/kong.conf
ENV KONG_PLUGINS bundled,custom-rate-limiting
ENV KONG_PROXY_ACCESS_LOG=/dev/stdout
ENV KONG_ADMIN_ACCESS_LOG=/dev/stdout
ENV KONG_PROXY_ERROR_LOG=/dev/stderr
ENV KONG_ADMIN_ERROR_LOG=/dev/stderr
ENV KONG_ADMIN_LISTEN="0.0.0.0:8001, 0.0.0.0:8444 ssl"
ENV KONG_NGINX_USER="root root"
ENV KONG_PROXY_LISTEN 0.0.0.0:8000
ENV KONG_LOG_LEVEL debug
USER root

RUN  mkdir -p /etc/kong/plugins
COPY --from=builder  /tmp/go/src/rate-limiting/go-pluginserver/go-pluginserver /usr/local/bin
COPY --from=builder  /tmp/go/src/rate-limiting/custom-rate-limiting.so /etc/kong/plugins

# for debug
RUN /usr/local/bin/go-pluginserver -version && \
    cd /etc/kong/plugins && \
    /usr/local/bin/go-pluginserver -dump-plugin-info custom-rate-limiting