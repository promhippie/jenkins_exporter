#!/bin/sh
set -e

if [ ! -d /var/lib/jenkins-exporter ]; then
    userdel jenkins-exporter 2>/dev/null || true
    groupdel jenkins-exporter 2>/dev/null || true
fi
