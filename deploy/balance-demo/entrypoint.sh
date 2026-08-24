#!/bin/sh
set -eu

instance="${DEMO_NAME:-demostracion}"
sed -e "s/__INSTANCE__/${instance}/g" -e "s/__CONTAINER_ID__/${HOSTNAME}/g" \
  /opt/balance/index.html.template > /usr/share/nginx/html/index.html

exec nginx -g 'daemon off;'
