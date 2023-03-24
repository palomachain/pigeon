#!/usr/bin/env bash
set -euo pipefail

function build_abigen_binary() {
  local project_dir="$1"
  local dependency='github.com/ethereum/go-ethereum'
  # getting the folder, where the go-ethereum package of the version, specified in go.mod, is located.
  local goeth_path=$(go list -m -f '{{.Dir}}' "$dependency")
  local temp_build_folder="/tmp/temp-goeth-build"

  if [[ -f "$project_dir/build/abigen" ]]; then
    while true; do
      read -p "Abigen binary already exists in the build directory. Do you want to rebuild it (y/n)?" yn
      case $yn in
        [Yy]* ) echo "Rebuilding ..."; break;;
        [Nn]* ) exit;;
        * ) echo "Please answer y or n.";;
      esac
    done
  fi

  if [[ ! -d "$goeth_path" ]]; then
    echo "No path found for the dependency $dependency, please run 'go mod download' or 'go mod tidy' first"
    exit 1
  fi

  rm -rf "$temp_build_folder"
  cp -rf "$goeth_path" "$temp_build_folder"
  chmod -R +w "$temp_build_folder"
  pushd "$temp_build_folder"
  make all
  mkdir -p "$project_dir/build/"
  cp "$temp_build_folder/build/bin/abigen" "$project_dir/build/"
  popd
  rm -rf "$temp_build_folder"
}

function abigen_generate_compass() {
  local project_dir="$1"
  local compass_abi_path="chain/evm/abi/compass"

  "$project_dir/build/abigen" --abi "$project_dir/$compass_abi_path/compass.abi" --pkg compass --out "$project_dir/$compass_abi_path/compass.go"
}
