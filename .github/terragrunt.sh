#!/bin/bash

# Safe mode
set -euo pipefail
IFS=$'\n\t'

ENV=$1

if [ -z "$ENV" ]; then
  echo "ENV variable is empty!\n"
  exit 1
fi

terragrunt --version

cd infra-live/non-prod/us-east-1/$ENV

terragrunt plan-all