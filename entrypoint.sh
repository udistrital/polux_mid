#!/usr/bin/env bash

set -e
set -u
set -o pipefail

if [ -n "${PARAMETER_STORE:-}" ]; then
  export POLUX_MID__DB_USER="$(aws ssm get-parameter --name /${PARAMETER_STORE}/polux_mid/db/username --output text --query Parameter.Value)"
  export POLUX_MID__DB_PASS="$(aws ssm get-parameter --with-decryption --name /${PARAMETER_STORE}/polux_mid/db/password --output text --query Parameter.Value)"
fi

exec ./main "$@"