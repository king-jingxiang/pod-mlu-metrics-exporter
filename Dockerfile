FROM debian:stretch-slim

ADD src/pod-mlu-metrics-exporter /usr/bin/pod-mlu-metrics-exporter

ENV NVIDIA_VISIBLE_DEVICES=all
ENV NVIDIA_DRIVER_CAPABILITIES=utility

ENTRYPOINT ["pod-mlu-metrics-exporter", "-logtostderr", "-v", "8"]
#ENTRYPOINT ["pod-mlu-metrics-exporter"]
