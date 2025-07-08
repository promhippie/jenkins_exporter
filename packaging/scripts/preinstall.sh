#!/bin/sh
set -e

if ! getent group jenkins-exporter >/dev/null 2>&1; then
    groupadd --system jenkins-exporter
fi

if ! getent passwd jenkins-exporter >/dev/null 2>&1; then
    useradd --system --create-home --home-dir /var/lib/jenkins-exporter --shell /bin/bash -g jenkins-exporter jenkins-exporter
fi
