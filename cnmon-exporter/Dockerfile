FROM ruanxingbaozi/cambricon-cnmon-base

ARG CNMON_VERSION=0.1

COPY cnmon-exporter /usr/local/bin

RUN chmod +x /usr/local/bin/cnmon-exporter

VOLUME /run/prometheus

ENTRYPOINT [ "cnmon-exporter" ]
