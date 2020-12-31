#!/bin/bash

# Safe mode
set -euo pipefail
IFS=$'\n\t'

ENV=$1
TERRAGRUNT_LIVE_DIR=""

# If prod, no sub directory exists. If non-prod, it does.
if [ $ENV == "prod" ]; then
    TERRAGRUNT_LIVE_DIR="infra-live/$ENV/us-east-1"
else
    TERRAGRUNT_LIVE_DIR="infra-live/non-prod/us-east-1/$ENV"
fi

terragrunt --version

cd $TERRAGRUNT_LIVE_DIR

terragrunt plan-all