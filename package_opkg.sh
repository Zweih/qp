#!/usr/bin/env bash
set -euo pipefail

VERSION="$1"
ARCH="$2"
BIN="$3"
NEWS_PATH="$4"
MANPAGE_PATH="$5"
OUTDIR="$6"

PKG_NAME="qp"

declare -A OPKG_ARCH_MAP
OPKG_ARCH_MAP[x86_64]="x86_64"
OPKG_ARCH_MAP[aarch64]="aarch64_generic aarch64_cortex-a53 aarch64_cortex-a72"
OPKG_ARCH_MAP[armv7h]="arm_cortex-a7 arm_cortex-a7_neon-vfpv4 arm_cortex-a9 arm_cortex-a9_vfpv3-d16 arm_cortex-a8_vfpv3 arm_cortex-a15_neon-vfpv4 arm_cortex-a5_vfpv4"
OPKG_ARCH_MAP[mipsle]="mipsel_24kc mipsel_74kc mipsel_24kf mipsel_74kf mipsel_mips32"

OPKGARCHES="${OPKG_ARCH_MAP[$ARCH]:-}"
if [[ -z "$OPKGARCHES" ]]; then
  echo "No OpenWrt arch mappings defined for $ARCH"
  exit 1
fi

for OPKGARCH in $OPKGARCHES; do
  echo "Packaging $ARCH → $OPKGARCH"

  PKGDIR="$OUTDIR/opkg/${PKG_NAME}_${VERSION}_${OPKGARCH}"
  CONTROL_DIR="$PKGDIR/CONTROL"

  mkdir -p "$CONTROL_DIR"
  mkdir -p "$PKGDIR/usr/bin"
  mkdir -p "$PKGDIR/usr/share/man/man1"
  mkdir -p "$PKGDIR/usr/share/doc/$PKG_NAME"

  install -m 755 "$BIN" "$PKGDIR/usr/bin/qp"
  install -m 644 "$MANPAGE_PATH" "$PKGDIR/usr/share/man/man1/qp.1"
  install -m 644 "$NEWS_PATH" "$PKGDIR/usr/share/doc/$PKG_NAME/NEWS"

  SIZE_KB=$(du -sk --exclude=CONTROL "$PKGDIR" | cut -f1)

  cat >"$CONTROL_DIR/control" <<EOF
Package: $PKG_NAME
Version: $VERSION
Architecture: $OPKGARCH
Maintainer: Fernando Nunez <me@fernandonunez.io>
Installed-Size: $SIZE_KB
Section: utils
License: GPL-3.0-only
Description: qp - Query Packages. A CLI tool for querying installed packages.
Homepage: https://github.com/Zweih/qp
EOF

  fakeroot sh -c "cd \"$OUTDIR\" && opkg-build -o 0 -g 0 \"${PKGDIR#$OUTDIR/}\" ."

  mv "$OUTDIR/${PKG_NAME}_${VERSION}_${OPKGARCH}.ipk" \
    "$OUTDIR/qp-v${VERSION}-${OPKGARCH}.ipk"

done
