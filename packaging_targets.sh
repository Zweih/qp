#!/usr/bin/env bash

deb_archs=(x86_64 armv7h aarch64)
brew_archs=(darwin_amd64 darwin_arm64)
aur_archs=(x86_64 armv7h aarch64)
opkg_archs=(x86_64 armv7h aarch64 mipsle)

should_package() {
  local target="$1"
  local arch="$2"
  local -n list="$target"

  for a in "${list[@]}"; do
    [[ "$a" == "$arch" ]] && return 0
  done

  return 1
}
