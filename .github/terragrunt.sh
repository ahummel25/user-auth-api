#!/bin/bash

# Safe mode
set -euo pipefail
IFS=$'\n\t'

terragrunt --version

cd infra-live/non-prod/us-east-1/dev

terragrunt plan-all