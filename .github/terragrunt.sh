#!/bin/bash

# Safe mode
set -euo pipefail
IFS=$'\n\t'

ENV=$1
TERRAGRUNT_LIVE_DIRS=()

# If prod, no sub directory exists. If non-prod, it does.
if [ $ENV == "prod" ]; then
    TERRAGRUNT_LIVE_DIRS=("infra/aws-infra-live/$ENV/us-east-1" "infra/mongo-infra-live/us-east-1/$ENV")
else
    TERRAGRUNT_LIVE_DIRS=("infra/aws-infra-live/non-prod/us-east-1/$ENV" "infra/mongo-infra-live/non-prod/us-east-1/$ENV")
fi

terragrunt --version

for TERRAGRUNT_LIVE_DIR in ${TERRAGRUNT_LIVE_DIRS[@]}; do
    cd $GITHUB_WORKSPACE/$TERRAGRUNT_LIVE_DIR

    module_dirs=($(ls -d */))

    for module in ${module_dirs[@]}; do
        current_dir=$GITHUB_WORKSPACE/$TERRAGRUNT_LIVE_DIR/$module
        cd $current_dir
        echo -e "Running terragrunt plan for $current_dir\n"

        terragrunt plan
        echo -e "\n"
        
        #terragrunt apply --auto-approve
        #echo -e "\n"
    done;
done;
