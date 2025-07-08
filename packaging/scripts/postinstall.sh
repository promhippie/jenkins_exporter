#!/bin/sh
set -e

chown -R jenkins-exporter:jenkins-exporter /var/lib/jenkins-exporter
chmod 750 /var/lib/jenkins-exporter

if [ -d /run/systemd/system ]; then
    systemctl daemon-reload

    if systemctl is-enabled --quiet jenkins-exporter.service; then
        systemctl restart jenkins-exporter.service
    fi
fi
