#!/usr/bin/env bash
set -euo pipefail

VERSION="$1"
ARCH="${2:-}"
BIN="$3"
NEWS_PATH="$4"
MANPAGE_PATH="$5"
COPYRIGHT_PATH="$6"
OUTDIR="$7"

if [[ -z "$ARCH" ]]; then
  name=$(basename "$BIN")
  if [[ "$name" =~ qp-(.+) ]]; then
    ARCH="${BASH_REMATCH[1]}"
  else
    echo "ERROR: Could not determine architecture from '$BIN'"
    exit 1
  fi
fi

case "$ARCH" in
x86_64) DEBARCH="amd64" ;;
aarch64) DEBARCH="arm64" ;;
armv7h) DEBARCH="armhf" ;;
*) echo "Unknown arch: $ARCH" && exit 1 ;;
esac

PKGDIR="$OUTDIR/deb/qp-v${VERSION}-${DEBARCH}"
mkdir -p "$PKGDIR/DEBIAN"
mkdir -p "$PKGDIR/usr/bin"
mkdir -p "$PKGDIR/usr/share/doc/qp"
mkdir -p "$PKGDIR/usr/share/man/man1"
SIZE_KB=$(du -sk --exclude=DEBIAN "$PKGDIR" | cut -f1)

cat >"$PKGDIR/DEBIAN/control" <<EOF
Package: qp
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${DEBARCH}
Installed-Size: ${SIZE_KB}
Maintainer: Fernando Nunez <me@fernandonunez.io>
Description: qp - Query Packages. A CLI tool for querying installed packages.
Homepage: https://github.com/Zweih/qp
EOF

install -m 755 "$BIN" "$PKGDIR/usr/bin/qp"
install -m 644 "$MANPAGE_PATH" "$PKGDIR/usr/share/man/man1/qp.1"
install -m 644 "$NEWS_PATH" "$PKGDIR/usr/share/doc/qp/NEWS"
install -m 644 "$COPYRIGHT_PATH" "$PKGDIR/usr/share/doc/qp/copyright"

dpkg-deb --build "$PKGDIR"
mv "$PKGDIR.deb" "$OUTDIR/qp-v${VERSION}-${DEBARCH}.deb"
