# qp - query packages

**qp** is a command-line program for Linux and macOS to query installed packages across ecosystems.

**qp** queries over 6x faster than native package searching while returning more comprehensive metadata than native package search solutions.

Query packages from `brew`, `pacman`, `apt`, `flatpak`, `snap`, `npm`, `pipx`, `dnf`, and `opkg`. Ecosystems are added frequently!

**qp** supports querying with full boolean logic for package metadata, dependency relations, and more.

You can find installation instructions [here](#Installation).

Check [features](#Features) to find out more.

Check [usage](#Usage) for all available commands + options.

![downloads-badge](https://zweih.github.io/repulse-analytics/downloads_badge.svg) ![clones-badge](https://zweih.github.io/repulse-analytics/clones_badge.svg)

<img src="https://zweih.github.io/repulse-analytics/qp-logo-query-packages.svg" alt="qp logo - query packages CLI tool logo">

[![AUR version - qp](https://img.shields.io/aur/version/qp?style=flat-square&logo=arch-linux&logoColor=1793d1&label=qp&color=1793d1)](https://aur.archlinux.org/packages/qp)
[![AUR version - qp-bin](https://img.shields.io/aur/version/qp-bin?style=flat-square&logo=arch-linux&logoColor=1793d1&label=qp-bin&color=1793d1)](https://aur.archlinux.org/packages/qp-bin)
[![AUR version - qp-git](https://img.shields.io/aur/version/qp-git?style=flat-square&logo=arch-linux&logoColor=1793d1&label=qp-git&color=1793d1)](https://aur.archlinux.org/packages/qp-git)

![Alt](https://repobeats.axiom.co/api/embed/a13406d103a649d70641774ee85e7a9983ccf96b.svg "Repobeats analytics image")

<details open>
<summary><strong>Download and clone statistics</strong></summary>
<br>
 
Graphs are generated daily with my other project, [Repulse Analytics](https://github.com/Zweih/repulse-analytics)

![repulse graphs for qp](https://zweih.github.io/repulse-analytics/combined_graphs.svg)

</details>

This package is compatible with the following platforms and distributions:
 - [Arch Linux](https://archlinux.org)
 - [macOS](https://www.apple.com/macos/)
 - [Debian](https://debian.org)
 - [SteamOS](https://store.steampowered.com/steamos)
 - [Red Hat Enterprise Linux (RHEL)](https://www.redhat.com/en/technologies/linux-platforms/enterprise-linux)
 - [Fedora](https://fedoraproject.org/)
 - [Linux Mint](https://linuxmint.com)
 - [Manjaro](https://manjaro.org/)
 - [Ubuntu](https://ubuntu.com)
 - [Pop!_OS](https://system76.com/pop/)
 - [CachyOS](https://cachyos.org/)
 - [OpenWrt](https://openwrt.org/)
 - [Garuda Linux](https://garudalinux.org/)
 - [EndeavourOS](https://endeavouros.com/)
 - [Mabox Linux](https://maboxlinux.org/)
 - [Zorin OS](https://zorin.com/os/)
 - [Elementary OS](https://elementary.io/)
 - The 50 other Arch, Debian, and Fedora-based distros, as long as they have `apt`/`dpkg`, `brew`, `pacman`, `flatpak`, `dnf`/`yum`, or `opkg` installed.

**qp** also detects and queries other system level package managers like `flatpak`, `npm`, and `pipx` for globally installed applications, expanding package discovery beyond traditional system package management.

**qp** supports embedded linux systems, including meta-distributions like [yocto](https://www.yoctoproject.org/) that use `opkg` (`.ipk` packages) or `apt`/`dpkg` (`.deb` packages) or `.rpm` packages.

## Features

* List installed packages across supported systems
* Compatible with MacOS, Arch, Debian, OpenWrt, and over 60 distros
  * Supports multiple ecosystems:
    * System package managers:
      * pacman, brew, apt/dpkg, dnf/yum, opkg
    * Application package managers:
      * Flatpak, npm, pipx
* Query packages using an expressive query language
  * Full boolean logic (`and`, `or`, `not`, grouping)
  * Fuzzy and strict matching
  * Range queries for `size`, `updated`, and `built`
  * Existence checks (`has:`, `no:`)
  * Learn more about querying [here](#Querying-with-where)
  * Complex queries via grouping (`q ... p`) and built-in macros
    * Includes `orphan` and `superorphan` filters
* Sort results by any field
* Output formats:
  * Table (default)
  * Key/value 
  * JSON
* Query by:
  * Name, version, origin, architecture, license
  * Cross-origin package detection (also-in field shows where packages exist across different package managers)
  * Size on disk, freeable storage, and total footprint
  * Update or build time/date
  * Package base or groups
  * Dependencies, optional dependencies, reverse dependencies
  * Package provisions, conflicts, replacements
  * Installation reason (explicit or dependency)
  * Package validation method (e.g., sha256, pgp)
  * Packager, URL, description
  * Package type (debug, split, cask, formula, etc.)
* Customizable field selection for output
* Cache system for fast repeated queries
* Lightweight, fast, concurrent architecture
* CLI designed for both scripting and interactive use
* Extensive roadmap with frequent improvements and optimizations

Learn about usage [here](#Usage).

learn about installation [here](#Installation).

## Is it good?

[Yes.](https://news.ycombinator.com/item?id=3067434)

## Roadmap

<details>
<summary><strong>Phase 1</strong></summary>


| Status | Feature | Status | Feature |
|--------|---------|--------|---------|
| ✓ | remove expac as a dependency (300% speed boost) | ✓ | concurrent file reading (200% speed boost) |
| ✓ | protobuf caching (127% speed boost) | ✓ | use chunked channel-based concurrent querying (12% speed boost) |
| ✓ | optimize file reading (28% speed boost) | ✓ | improve sorting efficiency (8% speed boost) |
| ✓ | optimize query order (4% speed boost) | ✓ | concurrent querying |
| ✓ | concurrent sorting | ✓ | asynchronous progress bar |
| ✓ | channel-based aggregation | ✓ | rewrite in golang |
| ✓ | automate binaries packaging | ✓ | CI to release binaries |
| ✓ | dependency depth resolution | ✓ | config dependency injection for testing |
| ✓ | query by package name | ✓ | query by size on disk |
| ✓ | query by range of size on disk | ✓ | query by date range |
| ✓ | user defined field selection | ✓ | dependencies of each package (dependency field) |
| ✓ | reverse-dependencies of each package (required-by field) | ✓ | list of packages for package queries |
| ✓ | package provisions | ✓ | package description query |
| ✓ | package conflicts field | ✓ | package architecture field |
| ✓ | package URL field | ✓ | package version field |
| ✓ | package license field | ✓ | package base field |
| ✓ | package base sort | ✓ | license query |
| ✓ | license sort | ✓ | dependency graph |
| ✓ | metaflag for all queries | ✓ | JSON output |
| ✓ | no-headers option | ✓ | provides query |
| ✓ | depends query | ✓ | all-fields option |
| ✓ | required-by query | ✓ | no cache option |
| ✓ | optional full timestamp | ✓ | package description field |
| ✓ | groups query | ✓ | driver interface |
| ✓ | package base query | ✓ | required-by sort |
| ✓ | optdepends sort | ✓ | depends sort |
| ✓ | build-date field | ✓ | build-date query |
| ✓ | build-date sort | ✓ | pkgtype field |
| ✓ | url query | ✓ | pkgtype sort |
| ✓ | architecture query | ✓ | groups field |
| ✓	| conflicts query | ✓ | package description sort |
| ✓	| regenerate cache option | ✓ | validation query |
| ✓ | url sort | ✓ | groups sort |
| ✓ | packager field | ✓ | optional dependency field |
| ✓ | sort by size on disk | ✓ | conflicts sort |
| ✓ | optional-for sort | ✓ | provides sort |
| ✓ | validation field | ✓ | validation sort |
| ✓ | packager sort | ✓ | architecture sort |
| ✓ | reason sort | ✓ | version sort |
| ✓ | pkgtype query | ✓ | optdepends query |
| ✓ | origin sort | ✓ | origin query |
| ✓ | packager query | ✓ | origin field |
| ✓ | replaces sort | ✓ | optional-for query 

</details>

<strong>Phase 2</strong>

| Status | Feature | Status | Feature |
|--------|---------|--------|---------|
| ✓ | reverse optional dependencies field (optional-for) | - | optdepends installation indicator |
| - | separate field for optdepends reason | ✓ | fuzzy/strict querying |
| ✓ | existence querying | ✓ | depth querying |
| ✓ | command-based syntax | ✓ | full boolean logic |
| ✓ | abstract syntax tree | ✓ | directed acyclical graph for filtering |
| - | user-defined macros | ✓ | parentetical (grouping) logic |
| ✓ | limit from end | ✓ | limit from middle |
| ✓ | built-in macros | – | streaming pipeline |
| - | query explaination | - | user configuration file |
| ✓ | deb origin (apt/dpkg support) | ✓ | deb packaging |
| ✓ | opkg origin (openwrt support) | ✓ | brew origin (homebrew support)|
| ✓ | bottles in brew | ✓ | casks in brew |
| - | replaced-by resolution | - | multi-license support |
| – | short-args for queries | ✓ | key/value output |
| ✓ | rpm origin (dnf/yum support) | ✓ | homebrew packaging |
| ✓ | pipx origin (python global packages) | ✓ | formulae from taps (brew) |
| - | casks from taps (brew) | - | dependencies for casks |
| - | rpm packaging | - | zypper (openSUSE support) |
| ✓ | cache-only option | ✓ | pacman hook |
| - | brew hook | - | deb hook |
| ✓ | npm origin (npm global packages) | - | nested dependencies |
| ✓ | also-in field (cross-origin managed) | ✓ | env field |
| ✓ | other envs field | - | support for multiple virtual environments (nvm/pyenv/etc.) |
| ✓ | freeable field | ✓ | footprint field |
| ✓ | flatpak origin | - | install history |
| ✓ | snap origin | ✓ | title field |
| - | keywords/tags field | - | notes/comment field |
| - | author field | - | cargo origin |
| - | log levels | ✓ | chunked cache (70% speed boost) |

## Installation

### Homebrew (macOS or Linuxbrew)

If you have [homebrew](https://brew.sh/) (`brew`), install via **qp**'s cask repo:
```bash
brew tap zweih/qp
brew install zweih/qp/qp
```

**Note**: Until we are added to the official `homebrew/core` repo, ensure that the package you install is `zweih/qp/qp`

### Arch-based systems (AUR)

Install using an [AUR helper](https://wiki.archlinux.org/title/AUR_helpers) like `yay`, `paru`, `aura`, etc.:
```bash
yay -Sy qp
```

If you prefer to install a pre-compiled binary* using the AUR, use the `qp-bin` package instead.

***Note**: Binaries are automatically, securely, and transparently compiled with github CI when a version release is created. You can audit the binary creation by checking the relevant github action for each release version.

For the latest (unstable) version from git w/ the AUR, use `qp-git`*.  

***Note**: `qp-git` is not recommended for most users

### Debian-based systems (e.g. Ubuntu, Mint, Pop!_os)

To install the latest `.deb` release for your system architecture:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/zweih/qp/packaging/install-qp-deb.sh)
```

This script downloads the latest published release from github and installs it using `dpkg`.

\***Note**: You can inspect the script beforehand [here](https://github.com/zweih/qp/packaging/install-qp-deb.sh). All binaries are built and signed by github CI on release tags.

Installation via `apt` is coming soon™!

### Building from source + manual installation

1. Clone the repo:
   ```bash
   git clone https://github.com/zweih/qp.git
   cd qp
   ```
2. Build the binary:
   ```bash
   go build -o qp ./cmd/qp
   ```
3. Copy the binary to your system's `$PATH`:
   ```bash
   sudo install -m755 qp /usr/bin/qp
   ```
4. Copy the manpage:
   ```bash
   sudo install -m644 qp.1 /usr/share/man/man1/qp.1
   ```

### Cache

#### Linux:
The cache is located under `/query-packages` at `$HOME/.cache/` or wherever you have `$XDG_HOME_CACHE` set to.

#### macOS:
The cache is located under `/query-packages` at `$HOME/Library/Caches/`

## Usage

```bash
qp [command] [args] [options]
```

### Commands

- `select <list>` | `s <list>`: Comma-separated list of fields to display
  - `select all` | `s all` will act as a list of all available fields
  - `select default` | `s default` will act as a list of all default fields
  - Use `select default,version` to list default fields + version
  - Use `select all,version` to list default fields + version
  - [See fields available for selection](#Available-Fields)
- `where <query>` | `w <query>`: Apply one or more queries to refine package results
  - Learn more about querying [here](#Querying-with-where)
- `order <field>:<direction>` | `o <field>:<direction>`: Sort results ascending or descending
  - Default sort is `updated:asc`
- `limit <number>` | `l <number>`: Limit the amount of packages to display (default: 20)
  - `limit all` | `l all`: Display all packages
  - `limit end:<number>`: Display last `n` packages
  - `limit mid:<number>`: Display middle `n` packages
- `format <type>` | `f <type>`: Output format: table, json, kv (default: table)
  - `format table` -> Tabular output with headers
  - `format json` -> JSON array output
  - `format kv` -> Key-Value pairs (best for selecting many fields)

### Options

- `--no-headers`: Omit column headers in table output (useful for scripting)
- `--full-timestamp`: Display the full timestamp (date and time) of package update/build instead of just the date
- `--no-cache`: Disable cache loading/saving and force fresh package data loading
- `--regen-cache`: Disable cache loading, force fresh package data loading, and save fresh cache
- `--cache-only`: Update cache only and nothing else. specify origin ('pacman', 'brew', 'deb', etc.) or all.
- `-h` | `--help`: Print help info

### Available fields

- `updated` - When the package was last updated
- `built` - When the package was built
- `size` - Package size on disk
- `freeable` - The amount storage that will be freed if the package is uninstalled
- `footprint` - The real-world proportional disk space impact of the package, accounting for shared dependencies
- `name` - Package name
- `title` - Alternate, full, or application name; some origins have a proper name for a package other than the name used to install it (i.e., brew has `vlc` and the title is `VLC Media Player`)
- `reason` - Installation reason (explicit/dependency)
- `version` - Installed package version
- `origin` - The package ecosystem or source the package belongs to (e.g., brew, pipx); reflects which package manager or backend maintains it
- `arch` - Architecture the package was built for (e.g., x86_64, aarch64, any)
- `env` - The environment where the package is installed (applies to origins that can have multiple virtual environments, such as npm)
- `license` - Package software license
- `description` - Package description
- `url` - The URL of the official site of the software being packaged
- `validation` - Package integrity validation method (e.g., sha256, pgp)
- `pkgtype` - Package type (specific to each origin, some origins have no pkgtype)
- `pkgbase` - Name of the base package used to group split packages; for non-split packages, it is the same as the package name. 
- `packager` - Person/entity who built the package (if available)
- `also-in` - List of other origins that have this package installed
- `other-envs` - List of other environments where this package is also installed
- `groups` - List of package groups or categories (e.g., base, gnome, xfce4)
- `conflicts` - List of packages that conflict, or cause problems, with the package
- `replaces` - List of packages that are replaced by the package
- `depends` - List of dependencies
- `optdepends` - List of optional dependencies
- `required-by` - List of packages required by the package and are dependent
- `optional-for` - List of packages that optionally depend on the package (optionally dependent)
- `provides` - List of alternative package names or shared libraries provided by package

### Package Types by Origin

The `pkgtype` field indicates the type or category of package within each ecosystem. Different origins use different package type classifications:

| Origin | Supported Package Types | Description |
|--------|------------------------|-------------|
| pacman | `pkg`, `split`, `debug`, `src` | Package build type |
| brew | `formula`, `cask` | Formulae are command-line tools, casks are GUI applications |
| flatpak | `app`, `runtime` | Runtimes are the main dependencies of apps |
| deb | none | - |
| rpm | none | - |
| opkg | none | - |
| pipx | none | - |
| npm | none | - |

**Notes:**
- Pacman's pkgtype comes from the package's XDATA field introduced in newer pacman versions
    - Older packages may not have this field populated
- Flatpak runtimes are always listed as dependencies
- Deb, rpm, opkg, and pipx origins do not implement package type classifications and will show empty pkgtype values

### Querying with `where`

The `where` (short: `w`) command is the core of **qp**'s flexible query system.

#### Logical Operations

- `and`: Combine multiple conditions where all must match
- `or`: Match any of multiple conditions
- `not`: Invert any condition
- `q ... p`: Qroup conditions for operation precedence, use them as you would parentheses with `q` being `(` and `p` being `)`
  - The purpose of this is that `(` and `)` are not safe to use unquoted on the command line
  - Remember with:
    - `q` is for **q**uery group start: `(`
    - `p` is for query group sto**p**: `)`

Examples:
```bash
qp w name=vim and size=10MB:
qp w reason=explicit and not required-by
qp w q name=zoxide or name=yazi p and optdepends=fzf
qp w q name=firefox or name=librewolf p and not has:conflicts
```

#### Query Types

All queries that take words as arguments can also take a comma-separated list.

Each `where` query supports one of the following:

- **String match**
  - `field=value` -> Fuzzy match
  - `field==value` -> Strict match
  - Applies to string-based fields like `name`, `license`, `description`, etc.
- **Range match**
  - `field=start:end` -> Fuzzy match
  - `field==start:end` -> Strict match
  - Works with `updated`, `built`, `size`, `freeable`, and `footprint`
  - Supports:
    - Full ranges: `start:end`
    - Open-ended ranges: `start:` or `:end`
    - Exact values: `value`
- **Existence check**
  - `has:field` -> Field must exist or be non-empty
  - `no:field` -> Field must be missing or empty

#### Match Types

|Field Type | Fuzzy | Strict | 
|-----------|-------|--------|
| Strings & Relations | Substring match (case-insensitive) | Exact character match (case-insensitive) |
| Dates | Matches by day (ignores time) | Matches exact timestamp (to the second) |
| Size  | ±0.3% byte tolerance (approximate) | Matches exact byte size |

For example:
  - `name=gtk` matches `gtk3`, `libgtk`, etc. (fuzzy)
  - `name==gtk` only matches a package named exactly `gtk`

#### Depth Querying

For relation fields, you can specify the depth level to query using the `@` syntax:

```bash
qp w depends=glibc@2    # packages that depend on glibc at depth 2
qp w required-by=gtk3@1 # packages directly required by gtk3 (depth 1)
qp w provides=libssl@3  # packages that provide libssl at depth 3
```

**Depth levels:**
- No `@` specified: depth 1 (direct relations only)
- `@1`: Direct relations only (same as default)
- `@2`: Second-level relations (relations of relations)
- `@3`, `@4`, etc.: Deeper levels in the dependency tree

**Specific for optional dependencies:**
- `optdepends` and `optional-for` return the optional relationships at the depth 1
- After depth 1, the dependency resolution includes hard dependencies from those optional packages in the final results. This is intentional.

**Examples:**
```bash
# show packages with direct dependencies on python (depth 1 implied)
qp w depends=python

# show packages that indirectly depend on openssl at depth 2
qp w depends=openssl@2
```

**Note:** Depth querying works with all relation fields: `depends`, `optdepends`, `required-by`, `optional-for`, `provides`, `conflicts`, and `replaces`.

#### Built-in Macros

Some frequently-used query patterns are available as built-in macros for convenience.

* `orphan` - Matches orphaned packages (dependencies no longer required by anything):
  ```
  qp where orphan
  ```

  is equivalent to:
  ```
  qp where no:required-by and reason=dependency
  ```

* `superorphan` - Matches "super" orphaned packages (dependencies no longer required by anything AND optional for nothing)
  ```
  qp w superorphan
  ```

  is equivalent to:
  ```
  qp where no:required-by and reason=dependency and no:optional-for
  ```

* `heavy` - Matches packages 100MB and larger
  ```
  qp w superorphan
  ```

  is equivalent to:
  ```
  qp where size=100MB:
  ```

* `light` - Matches packages 1MB and smaller
  ```
  qp w light
  ```

  is equivalent to:
  ```
  qp where size=1MB:
  ```

These macros can be combined with other queries as usual:

```
qp w orphan and size=100KB:
qp w not superorphan and not name=gtk
```

#### Query Examples

```
qp w name==bash and has:depends
qp where size=100MB:2GB
qp w name=python,cmake,yazi
qp w updated==2024-01-01
qp where q name=vim or name=nvim p and not has:conflicts
qp w not arch=x86_64
qp w q has:depends or has:required-by p and not reason=explicit
```

#### Field Types

| Field Type | Description |
|------------|-------------|
| String | Matches textual fields. Used for fields like name, license, description, etc. <br> can take a comma-separated list |
| Range | Matches numerical or time-based fields across a range. <br> Supports full ranges (start:end), open-ended ranges (start: / :end), or exact values |
| Relation | Matches fields that contain relationships to other packages (e.g., dependencies, conflicts, provides) <br> Can take a comma-separated list |

#### Fields and Their Types

| Field Name | Field Type |
|------------|------------|
| updated | range |
| built | range |
| size | range |
| freeable | range |
| footprint | range |
| name | string |
| reason | string |
| version | string |
| origin | string |
| arch | string |
| env | string |
| license | string |
| pkgbase | string |
| description | string |
| url | string
| validation | string |
| pkgtype | string |
| packager | string |
| groups | string |
| also-in | string | 
| other-envs | string |
| conflicts | relation |
| replaces | relation |
| depends | relation |
| optdepends | relation |
| required-by | relation |
| optional-for | relation |
| provides | relation |

### JSON output

`format json` outputs the package data as structured JSON instead of a table. this can be useful for scripts or automation.

Example:
```
qp select all where name=gtk3 format json
```

Output format:
```json
[
  {
    "installTimestamp": 1743448253,
    "buildTimestamp": 1741400060,
    "size": 58266727,
    "freeable": 174966161,
    "footprint": 335629243,
    "name": "gtk3",
    "reason": "dependency",
    "version": "1:3.24.49-1",
    "origin": "pacman",
    "arch": "x86_64",
    "license": "LGPL-2.1-or-later",
    "description": "GObject-based multi-platform GUI toolkit",
    "url": "https://www.gtk.org/",
    "validation": "pgp",
    "pkgtype": "split",
    "pkgbase": "gtk3",
    "packager": "Jan Alexander Steffens (heftig) <heftig@archlinux.org>",
    "conflicts": [
      "gtk3-print-backends"
    ],
    "replaces": [
      "gtk3-print-backends<=3.22.26-1"
    ],
    "depends": [
      "adwaita-icon-theme",
      "at-spi2-core",
      "cairo",
      "cantarell-fonts",
      "dconf",
      "desktop-file-utils",
      "fontconfig",
      "fribidi",
      "gdk-pixbuf2",
      "glib2",
      "glibc",
      "gtk-update-icon-cache",
      "harfbuzz",
      "iso-codes",
      "libcloudproviders",
      "libcolord",
      "libcups",
      "libegl → libglvnd",
      "libepoxy",
      "libgl → libglvnd",
      "librsvg",
      "libx11",
      "libxcomposite",
      "libxcursor",
      "libxdamage",
      "libxext",
      "libxfixes",
      "libxi",
      "libxinerama",
      "libxkbcommon",
      "libxrandr",
      "libxrender",
      "pango",
      "shared-mime-info",
      "tinysparql",
      "wayland"
    ],
    "optDepends": [
      "evince (Default print preview command)"
    ],
    "requiredBy": [
      "ibus",
      "libdbusmenu-gtk3"
    ],
    "optionalFor": [
      "avahi (avahi-discover, avahi-discover-standalone, bshell, bssh, bvnc)",
      "libdecor (gtk3 support)",
      "pinentry (GTK backend)"
    ],
    "provides": [
      "gtk3-print-backends",
      "libgailutil-3.so=0-64",
      "libgdk-3.so=0-64",
      "libgtk-3.so=0-64"
    ]
  }
]
```

### Tips & Tricks

- Multiple short commands are supported using space separation (e.g. `s`, `w`, `l`, `o`, `f`), but **cannot** be combined as `swo` or `-swo`. use them like this:
  ```
  qp w name yay
  qp s name,size w name=vim o updated:asc l 10 # full query with shorthand
  ```

- Qroup queries with `q ... p` to clarify order of operations:
  ```
  qp w q name=git or name=gh p and has:depends 
  ```

- The relations table columns can be lengthy. Packages like `glibc` are required by thousands of packages. to improve readability, pipe the output to tools like `moar` or `less` (i prefer `moar`, but `less` is usually pre-installed):
  ```
  qp select name,depends | less
  qp s name,depends | moar
  ```

- Arguments to queries can be quoted if they contain special characters or spaces:
  ```
  qp where description="for tree-sitter"
  ```

- The `--no-headers` flag is useful when processing output in scripts. It removes the header row, making it easier to parse package lists with tools like `awk`, `sed`, or `cut`:
  ```
  qp --no-headers select name,size | awk '{print $1, $2}'
  ```

### Examples

 1. Show the last 10 installed packages  
   ```
   qp limit 10
   ```

 2. Show all explicitly installed packages  
   ```
   qp where reason=explicit limit all
   ```

 3. Show only dependencies updated on a specific date  
   ```
   qp where reason=dependency and update=2025-03-01
   ```

 4. Show all packages sorted alphabetically by name
   ```
   qp order name limit all
   ```

 5. Search for packages that contain a GPL license 
   ```
   qp where license=gpl
   ```

 6. Show packages installed between January 1, 2025, and January 5, 2025  
   ```
   qp where updated=2025-01-01:2025-01-05
   ```

 7. Sort all packages by their license, displaying name and license  
   ```
   qp select name,license order license limit all
   ```

 8. Show the 20 most recently installed packages larger than 20MB  
   ```
   qp where size=20MB: limit 20
   ```

 9. Show packages between 100MB and 1GB installed up to February 27, 2025
   ```
   qp where size=100MB:1GB and updated=:2025-02-27
   ```

10. Show all packages sorted by size in descending order, installed after January 1, 2025
   ```
   qp where updated=2025-01-01: order size:desc limit all
   ```

11. Search for installed packages containing "python"  
   ```
   qp where name=python
   ```

12. Search for explicitly installed packages containing "lib" that are between 10MB and 1GB in size
   ```
   qp where reason=explicit and name=lib and size=10MB:1GB
   ```

13. Search for packages with names containing "linux" installed between January 1 and March 30, 2025
   ```
   qp where name=linux and updated=2025-01-01:2025-03-30
   ```

14. Search for packages containing "gtk" installed after January 1, 2025, and at least 5MB in size
   ```
   qp where name=gtk and updated=2025-01-01: and size=5MB:
   ```

15. Show packages with name, version, and size
   ```
   qp select name,version,size
   ```

16. Show package names, descriptions, and dependencies with `less` for readability
   ```
   qp select name,depends,description | less
   ```

17. Output package data in JSON format
   ```
   qp format json
   ```

18. Save all explicitly installed packages to a JSON file
   ```
   qp where reason=explicit format json > explicit-packages.json
   ```

19. Output all packages sorted by size (descending) in JSON
   ```
   qp order size:desc limit all format json
   ```

20. Output JSON with specific fields
   ```
   qp select name,version,size format json
   ```

21. Show all available package details for all packages
   ```
   qp select all limit all
   ```

22. Output all packages with all fields in JSON format
   ```
   qp select all limit all format json
   ```

23. Show package names and sizes without headers for scripting
   ```
   qp select name,size --no-headers
   ```

24. Show all packages required by `firefox`
   ```
   qp where required-by=firefox limit all
   ```

25. Show all packages required by `gtk3` that are at least 50MB in size
   ```
   qp where required-by=gtk3 and size=50MB: limit all
   ```

26. Show packages required by `vlc` and installed after January 1, 2025
   ```
   qp where required-by=vlc and updated=2025-01-01:
   ```

27. Show all packages that have `glibc` as a dependency and are required by `ffmpeg`
   ```
   qp where depends=glibc and required-by=ffmpeg limit all
   ```

28. Inclusively show packages that require `gcc` or `pacman`
   ```
   qp where required-by=gcc,pacman
   ```

29. Show packages that provide `awk`
   ```
   qp where provides=awk
   ```

30. Inclusively show packages that provide `rustc` or `python3`
   ```
   qp where provides=rustc,python3
   ```

31. Show packages that conflict with `linuxqq`
   ```
   qp where conflicts=linuxqq
   ```

32. Show packages that are built for the `aarch64` CPU architecture or any architecture
   ```
   qp where arch=aarch64,any
   ```

33. Show all dependencies smaller than 500KB
   ```
   qp where reason=dependency and size=:500KB
   ```

34. Show the 15 most recent explicitly installed packages
   ```
   qp where reason=explicit limit 15
   ```

35. Show packages that contain "clang" in their description
   ```
   qp where description=clang
   ```

36. Sort packages by their package base while showing their names and package bases, in reverse alphabetical order
   ```
   qp select name,pkgbase order pkgbase:desc
   ```

37. Show packages that are exactly named "bash"
   ```
   qp where name==bash
   ```

38. Show packages that have no dependencies
   ```
   qp where no:depends
   ```

39. Show all packages installed via pipx
   ```
   qp where origin=pipx
   ```

## License

This project is licensed under GPL-3.0-only.

For use cases not compatible with the GPL, such as proprietary redistribution or integration/ingestion into ML/LLM systems, a separate commercial license is available. See LICENSE.commercial for details.
