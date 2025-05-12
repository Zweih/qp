#!/usr/bin/env bash
set -euo pipefail

RELEASE_DIR="./release"
EXTRA_DIR="./release-extra"

mkdir -p "$EXTRA_DIR"

CORE_ARCHES=(
  x86_64
  aarch64
  aarch64_generic
  arm_cortex-a9_vfpv3-d16
  mipsel_24k
)

CORE_ARCH_REGEX=$(
  IFS=\|
  echo "${CORE_ARCHES[*]}"
)

echo "Filtering .ipk packages in $RELEASE_DIR..."

shopt -s nullglob
for file in "$RELEASE_DIR"/qp-v*-*.ipk; do
  base=$(basename "$file")
  arch=$(echo "$base" | sed -n 's/^qp-v[0-9.]*-\(.*\)\.ipk$/\1/p')

  if [[ -z "$arch" ]]; then
    echo "Skipping malformed filename: $base"
    continue
  fi

  if [[ ! "$arch" =~ ^($CORE_ARCH_REGEX)$ ]]; then
    echo "Moving $base to release-extra/"
    mv "$file" "$EXTRA_DIR/"
  else
    echo "Keeping $base in main release"
  fi
done
