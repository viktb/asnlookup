#!/usr/bin/env bash
set -euo pipefail

base_url="http://archive.routeviews.org/bgpdata"
filename=$(curl -fsSL "${base_url}/$(date +%Y.%m)/RIBS/" | grep -oE 'rib\.[0-9]{8}\.[0-9]{4}\.bz2' | tail -n 1)
curl -fSLO "${base_url}/$(date +%Y.%m)/RIBS/${filename}"
