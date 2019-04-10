#!/usr/bin/env bash
DATE=`date +%Y-%m-%d-%H`
git add .
git commit -m "$DATE"
git push origin master