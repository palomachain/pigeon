#!/usr/bin/env bash

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

exec reflex \
  -d none \
  -s \
  -r '\.(go|proto|yaml)$' \
  -R '^third_party_proto/' \
  -- \
  $@
