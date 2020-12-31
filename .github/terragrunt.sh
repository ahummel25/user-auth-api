#!/bin/bash

# Safe mode
set -euo pipefail
IFS=$'\n\t'

ENV=$1
SUB_ENV=$2

terragrunt --version

cd infra-live/$ENV/us-east-1/$SUB_ENV

terragrunt plan-all