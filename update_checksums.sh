#!/bin/bash

PKGBUILD_FILE=${1:-PKGBUILD}
SHA_FILE=${2:-SHA256SUMS.txt}

PKG_NAME="qp"
EXT="tar.gz"

[[ ! -f "$SHA_FILE" ]] && {
  echo "Error: $SHA_FILE not found!"
  exit 1
}

[[ ! -f "$PKGBUILD_FILE" ]] && {
  echo "Error: $PKGBUILD_FILE not found!"
  exit 1
}

get_checksum() {
  grep "$1" "$SHA_FILE" | awk '{print $1}'
}

pkgver=$(grep '^pkgver=' "$PKGBUILD_FILE" | cut -d'=' -f2)
version="v${pkgver}"

for arch in x86_64 aarch64 armv7h; do
  bin_file="${PKG_NAME}-${version}-${arch}.${EXT}"
  checksum_arch=$(get_checksum "${bin_file}")

  if [[ -n "$checksum_arch" ]]; then
    sed -i "/^sha256sums_${arch}=.*$/c\sha256sums_${arch}=('$checksum_arch')" "${PKGBUILD_FILE}"
    echo "Updated checksum for ${bin_file}"
  else
    echo "Warning: Checksum for ${bin_file} not found in ${SHA_FILE}!"
  fi
done

source_file="${PKG_NAME}-${version}.${EXT}"
checksum_source=$(get_checksum "${source_file}")

if [[ -z "$pkgver" ]]; then
  echo "Error: pkgver not found in ${PKGBUILD_FILE}!"
  exit 1
fi

if [[ -n "$checksum_source" ]]; then
  sed -i "s|^sha256sums=.*|sha256sums=('$checksum_source')|" "${PKGBUILD_FILE}"
  echo "Updated source checksum for ${source_file}"
else
  echo "Warning: Checksum for ${source_file} not found in ${SHA_FILE}!"
fi

echo "All checksums updated successfully."
