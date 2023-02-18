#!/usr/bin/env bash
set -e

val=$(printf '%s\n' "${!1}")
if [ "${val}" == "" ]; then
    echo "Wrong env var name: '${1}'"
    exit 1
else
    echo ${val}
fi
