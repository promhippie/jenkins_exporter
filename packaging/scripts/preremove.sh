#!/bin/sh
set -e

systemctl stop jenkins-exporter.service || true
systemctl disable jenkins-exporter.service || true
