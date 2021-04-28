#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://archive.routeviews.org/bgpdata"
TIMESTAMP=$((($(date -u +%s)-3600)/7200*7200))
YEAR=$(date -u -d "@${TIMESTAMP}" +%Y)
MONTH=$(date -u -d "@${TIMESTAMP}" +%m)
FILENAME=$(curl -fsSL "${BASE_URL}/${YEAR}.${MONTH}/RIBS/" | grep -oE 'rib\.[0-9]{8}\.[0-9]{4}\.bz2' | tail -n 1)
wget "${BASE_URL}/${YEAR}.${MONTH}/RIBS/${FILENAME}"
