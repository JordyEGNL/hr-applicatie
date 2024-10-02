#!/usr/bin/env bash

while getopts 'v:' flag; do
  case "${flag}" in
  v) VERSION=${OPTARG} ;;
  *)
    echo "Invalid args"
    exit 1
    ;;
  esac
done

sed -i 's/Version:.*/Version: \"$VERSION\"' src/frontend/login/index.html
sed -i 's/TAG=.*/TAG=\"$VERSION\"' production.env