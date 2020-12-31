#!/bin/bash

# Safe mode
set -euo pipefail
IFS=$'\n\t'

ENV=$1

terragrunt --version

cd infra-live/non-prod/us-east-1/$ENV

terragrunt plan-all