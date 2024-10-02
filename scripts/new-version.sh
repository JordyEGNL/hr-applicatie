#!/usr/bin/env bash

# Opties verwerken
while getopts 'v:' flag; do
  case "${flag}" in
    v) VERSION=${OPTARG} ;;  # Stel de VERSION variabele in op de opgegeven waarde
    *)
      echo "Invalid args"
      exit 1
      ;;
  esac
done

# Controleer of de VERSION variabele is ingesteld
if [ -z "$VERSION" ]; then
  echo "Version not provided. Use -v to specify the version."
  exit 1
fi

# Update de versie in index.html
sed -i "s/\(Version: \)[^<]*/\1$VERSION/" src/frontend/login/index.html

# Update de TAG in production.env
sed -i "s/TAG=.*/TAG=$VERSION/" production.env