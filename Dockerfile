FROM quay.io/prometheus/busybox:latest
MAINTAINER Thomas Boerger <thomas@webhippie.de>

COPY jenkins_exporter /bin/jenkins_exporter

EXPOSE 9103
ENTRYPOINT ["/bin/jenkins_exporter"]
