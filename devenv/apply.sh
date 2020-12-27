#!/usr/bin/env sh

fly -t dev set-pipeline -p github-pr-comments-resource -c pipeline.yml -l vars.yml
