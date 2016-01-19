#!/bin/bash
set -e

if [ "${1:0:1}" = '-' ]; then
    set -- river_metrics_transformer "$@"
fi

exec "$@"
