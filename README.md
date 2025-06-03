# qp - query packages

**qp** is a command-line program for linux and macOS to query installed packages across ecosystems.

**qp** queries over 6x faster than native package searching while returning more comprehensive metadata than native package search solutions.

query packages from `brew`, `pacman`, `apt`/`dpkg`, `pipx`, and `dnf`/`yum`. ecosystems are added frequently!

**qp** supports querying with full boolean logic for package metadata, dependency relations, and more.

you can find installation instructions [here](#installation).

check [features](#features) to find out more.

check [usage](#usage) for all available commands + options.

![downloads-badge](https://zweih.github.io/repulse-analytics/downloads_badge.svg) ![clones-badge](https://zweih.github.io/repulse-analytics/clones_badge.svg)

<img src="https://zweih.github.io/repulse-analytics/qp-logo-query-packages.svg" alt="qp logo - query packages CLI tool logo">

[![AUR version - qp](https://img.shields.io/aur/version/qp?style=flat-square&logo=arch-linux&logoColor=1793d1&label=qp&color=1793d1)](https://aur.archlinux.org/packages/qp)
[![AUR version - qp-bin](https://img.shields.io/aur/version/qp-bin?style=flat-square&logo=arch-linux&logoColor=1793d1&label=qp-bin&color=1793d1)](https://aur.archlinux.org/packages/qp-bin)
[![AUR version - qp-git](https://img.shields.io/aur/version/qp-git?style=flat-square&logo=arch-linux&logoColor=1793d1&label=qp-git&color=1793d1)](https://aur.archlinux.org/packages/qp-git)

![Alt](https://repobeats.axiom.co/api/embed/a13406d103a649d70641774ee85e7a9983ccf96b.svg "Repobeats analytics image")

<details open>
<summary><strong>download and clone statistics</strong></summary>
<br>
 
graphs are generated daily with my other project, [repulse analytics](https://github.com/Zweih/repulse-analytics)

![repulse graphs for qp](https://zweih.github.io/repulse-analytics/combined_graphs.svg)

</details>

this package is compatible with the following platforms and distributions:
 - [arch linux](https://archlinux.org)
 - [macOS](https://www.apple.com/macos/)
 - [debian](https://debian.org)
 - [steamOS](https://store.steampowered.com/steamos)
 - [red hat enterprise linux (RHEL)](https://www.redhat.com/en/technologies/linux-platforms/enterprise-linux)
 - [fedora](https://fedoraproject.org/)
 - [linux mint](https://linuxmint.com)
 - [manjaro](https://manjaro.org/)
 - [ubuntu](https://ubuntu.com)
 - [pop!_OS](https://system76.com/pop/)
 - [cachyOS](https://cachyos.org/)
 - [openwrt](https://openwrt.org/)
 - [garuda linux](https://garudalinux.org/)
 - [endeavourOS](https://endeavouros.com/)
 - [mabox linux](https://maboxlinux.org/)
 - [zorin OS](https://zorin.com/os/)
 - [elementary OS](https://elementary.io/)
 - the 50 other arch, debian, and fedora-based distros, as long as they have `apt`/`dpkg`, `brew`, `pacman`, `dnf`/`yum`, or `opkg` installed

**qp** also detects and queries other system level package managers like `pipx` for isolated python applications, expanding package discovery beyond traditional system package management.

**qp** supports embedded linux systems, including meta-distributions like [yocto](https://www.yoctoproject.org/) that use `opkg` (`.ipk` packages) or `apt`/`dpkg` (`.deb` packages) or `.rpm` packages!

## features

* list installed packages across supported systems
* compatible with macOS, arch, debian, openwrt, and over 60 distros
* * supports multiple ecosystems:
  * system package managers:
    * pacman, brew, apt/dpkg, dnf/yum, opkg
* * application package managers:
    * pipx
* query packages using an expressive query language
  * supports full boolean logic (`and`, `or`, `not`, grouping)
  * supports fuzzy and strict matching
  * supports range queries for `size`, `updated`, and `built`
  * supports presence/absence checks (`has:`, `no:`)
  * learn more about querying [here](#querying-with-where)
* sort results by any field
* output formats:
  * table (default)
  * key/value
  * JSON
* query by:
  * name, version, origin, architecture, license
  * size on disk
  * update or build time/date
  * package base or groups
  * dependencies, optional dependencies, reverse dependencies
  * package provisions, conflicts, replacements
  * installation reason (explicit or dependency)
  * package validation method (e.g., sha256, pgp)
  * packager, URL, description
  * package type (debug, split, etc.)
* complex queries via grouping (`q ... p`) and built-in macros
  * includes `orphan` and `superorphan` filters
* customizable field selection for output
* cache system for fast repeated queries
* lightweight, fast, concurrent architecture
* CLI designed for both scripting and interactive use
* extensive roadmap with frequent improvements and optimizations

learn about usage [here](#usage)

learn about installation [here](#installation)

## is it good?

[yes.](https://news.ycombinator.com/item?id=3067434)

## roadmap

<details>
<summary><strong>phase 1</strong></summary>


| Status | Feature | Status | Feature |
|--------|---------|--------|---------|
| ✓ | remove expac as a dependency (300% speed boost) | ✓ | concurrent file reading (200% speed boost) |
| ✓ | protobuf caching (127% speed boost) | ✓ | use chunked channel-based concurrent querying (12% speed boost) |
| ✓ | optimize file reading (28% speed boost) | ✓ | improve sorting efficiency (8% speed boost) |
| ✓ | optimize query order (4% speed boost) | ✓ | concurrent querying |
| ✓ | concurrent sorting | ✓ | asynchronous progress bar |
| ✓ | channel-based aggregation | ✓ | rewrite in golang |
| ✓ | automate binaries packaging | ✓ | add CI to release binaries |
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

<strong>phase 2</strong>

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
| ✓ | pipx origin (python global packages) | - | formulae from taps (brew) |
| - | casks from taps (brew) | - | dependencies for casks |
| - | rpm packaging | - | zypper (openSUSE support) |
| ✓ | cache-only option | ✓ | pacman hook |
| - | brew hook | - | deb hook |

## installation

### homebrew (macOS or linuxbrew)

if you have [homebrew](https://brew.sh/) (`brew`), install via **qp**'s cask repo:
```bash
brew tap zweih/qp
brew install zweih/qp/qp
```

**note**: until we are added to the official `homebrew/core` repo, ensure that the package you install is `zweih/qp/qp`

### arch-based systems (AUR)

install using an [AUR helper](https://wiki.archlinux.org/title/AUR_helpers) like `yay`, `paru`, `aura`, etc.:
```bash
yay -Sy qp
```

if you prefer to install a pre-compiled binary* using the AUR, use the `qp-bin` package instead.

***note**: binaries are automatically, securely, and transparently compiled with github CI when a version release is created. you can audit the binary creation by checking the relevant github action for each release version.

for the latest (unstable) version from git w/ the AUR, use `qp-git`*.  

***note**: this is not recommended for most users

### debian-based systems (e.g. ubuntu, mint, pop!_os)

to install the latest `.deb` release for your system architecture:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/zweih/qp/packaging/install-qp-deb.sh)
```

this script downloads the latest published release from github and installs it using `dpkg`.

\***note**: you can inspect the script beforehand [here](https://github.com/zweih/qp/packaging/install-qp-deb.sh). all binaries are built and signed by github CI on release tags.

installation via `apt` is coming soon™!

### building from source + manual installation

**note**: this packages is specific to arch-based and debian-based linux distributions

1. clone the repo:
   ```bash
   git clone https://github.com/zweih/qp.git
   cd qp
   ```
2. build the binary:
   ```bash
   go build -o qp ./cmd/qp
   ```
3. copy the binary to your system's `$PATH`:
   ```bash
   sudo install -m755 qp /usr/bin/qp
   ```
4. copy the manpage:
   ```bash
   sudo install -m644 qp.1 /usr/share/man/man1/qp.1
   ```

### cache

#### linux:
the cache is located under `/query-packages` at `$HOME/.cache/` or wherever you have `$XDG_HOME_CACHE` set to.

#### macOS:
the cache is located under `/query-packages` at `$HOME/Library/Caches/`

## usage

```bash
qp [command] [args] [options]
```

### commands

- `select <list>` | `s <list>`: comma-separated list of fields to display
  - `select all` | `s all` will act as a list of all available fields
  - `select default` | `s default` will act as a list of all default fields
  - use `select default,version` to list default fields + version
  - use `select all,version` to list default fields + version
  - [see fields available for selection](#available-fields)
- `where <query>` | `w <query>`: apply one or more queries to refine package results.
  - learn more about querying [here](#querying-with-where)
- `order <field>:<direction>` | `o <field>:<direction>`: sort results ascending or descending
  - default sort is `updated:asc`
- `limit <number>` | `l <number>`: limit the amount of packages to display (default: 20)
  - `limit all` | `l all`: display all packages
  - `limit end:<number>`: display last `n` packages
  - `limit mid:<number>`: display middle `n` packages

### options

- `--no-headers`: omit column headers in table output (useful for scripting)
- `--full-timestamp`: display the full timestamp (date and time) of package update/build instead of just the date
- `--output`: format output as `table`, `kv` (key-value), `json` (default:`table`)
- `--no-progress`: force no progress bar outside of non-interactive environments
- `--no-cache`: disable cache loading/saving and force fresh package data loading
- `--regen-cache`: disable cache loading, force fresh package data loading, and save fresh cache
- `--cache-only`: update cache only and nothing else. specify origin ('pacman', 'brew', 'deb', etc.) or all.
- `-h` | `--help`: print help info

### available fields

- `updated` - when the package was last updated
- `built` - when the package was built
- `size` - package size on disk
- `name` - package name
- `reason` - installation reason (explicit/dependency)
- `version` - installed package version
- `origin` - the package ecosystem or source the package belongs to (e.g., brew, pipx); reflects which package manager or backend maintains it
- `arch` - architecture the package was built for (e.g., x86_64, aarch64, any)
- `license` - package software license
- `description` - package description
- `url` - the URL of the official site of the software being packaged
- `validation` - package integrity validation method (e.g., sha256, pgp)
- `pkgtype` - package type (specific to each origin, some origins have no pkgtype)
- `pkgbase` - name of the base package used to group split packages; for non-split packages, it is the same as the package name. 
- `packager` - person/entity who built the package (if available)
- `groups` - list of package groups or categories (e.g., base, gnome, xfce4)
- `conflicts` - list of packages that conflict, or cause problems, with the package
- `replaces` - list of packages that are replaced by the package
- `depends` - list of dependencies
- `optdepends` - list of optional dependencies
- `required-by` - list of packages required by the package and are dependent
- `optional-for` - list of packages that optionally depend on the package (optionally dependent)
- `provides` - list of alternative package names or shared libraries provided by package

### package types by origin

the `pkgtype` field indicates the type or category of package within each ecosystem. different origins use different package type classifications:

| origin | supported package types | description |
|--------|------------------------|-------------|
| pacman | `pkg`, `split`, `debug`, `src` | package build type |
| brew | `formula`, `cask` | formulae are command-line tools, casks are GUI applications |
| deb | none | - |
| rpm | none | - |
| opkg | none | - |
| pipx | none | - |

**notes:**
- pacman's pkgtype comes from the package's XDATA field introduced in newer pacman versions
    - older packages may not have this field populated
- brew distinguishes between formulae (CLI tools/libraries) and casks (GUI applications)
- deb, rpm, opkg, and pipx origins do not implement package type classifications and will show empty pkgtype values

### querying with `where`

the `where` (short: `w`) command is the core of qp's flexible query system.

#### logical operations

- `and`: combine multiple conditions where all must match
- `or`: match any of multiple conditions
- `not`: invert any condition
- `q ... p`: group conditions for operation precedence, use them as you would parentheses with `q` being `(` and `p` being `)`
    - the purpose of this is that `(` and `)` are not safe to use unquoted on the command line
    - remember with:
        - `q` is for **q**uery group start: `(`
        - `p` is for query group sto**p**: `)`

examples:
```bash
qp w name=vim and size=10MB:
qp w reason=explicit and not required-by
qp w q name=zoxide or name=yazi p and optdepends=fzf
qp w q name=firefox or name=librewolf p and not has:conflicts
```

#### query types

all queries that take words as arguments can also take a comma-separated list.

each `where` query supports one of the following:

- **string match**
  - `field=value` -> fuzzy match
  - `field==value` -> strict match
  - applies to string-based fields like `name`, `license`, `description`, etc.
- **range match**
  - `field=start:end` -> fuzzy match
  - `field==start:end` -> strict match
  - works with `updated`, `built`, and `size`
  - supports:
    - full ranges: `start:end`
    - open-ended ranges: `start:` or `:end`
    - exact values: `value`
- **existence check**
  - `has:field` -> field must exist or be non-empty
  - `no:field` -> field must be missing or empty

#### match types

|field type | fuzzy | strict | 
|-----------|-------|--------|
| strings & relations | substring match (case-insensitive) | exact character match (case-insensitive) |
| dates | matches by day (ignores time) | matches exact timestamp (to the second) |
| size  | ±0.3% byte tolerance (approximate) | matches exact byte size |

for example:
  - `name=gtk` matches `gtk3`, `libgtk`, etc. (fuzzy)
  - `name==gtk` only matches a package named exactly `gtk`

#### depth querying

for relation fields, you can specify the depth level to query using the `@` syntax:

```bash
qp w depends=glibc@2    # packages that depend on glibc at depth 2
qp w required-by=gtk3@1 # packages directly required by gtk3 (depth 1)
qp w provides=libssl@3  # packages that provide libssl at depth 3
```

**depth levels:**
- no `@` specified: depth 1 (direct relations only)
- `@1`: direct relations only (same as default)
- `@2`: second-level relations (relations of relations)
- `@3`, `@4`, etc.: deeper levels in the dependency tree

**spefor optional dependencies:**
- `optdepends` and `optional-for` return the optional relationships at the depth 1
- after epth 1, the dependency resolution includes hard dependencies from those optional packages in the final results. this is intentional.

**examples:**
```bash
# show packages with direct dependencies on python (depth 1 implied)
qp w depends=python

# show packages that indirectly depend on openssl at depth 2
qp w depends=openssl@2
```

**note:** depth querying works with all relation fields: `depends`, `optdepends`, `required-by`, `optional-for`, `provides`, `conflicts`, and `replaces`.

#### built-in macros

some frequently-used query patterns are available as built-in macros for convenience.

* `orphan` - matches orphaned packages (dependencies no longer required by anything):
  ```
  qp where orphan
  ```

  is equivalent to:
  ```
  qp where no:required-by and reason=dependency
  ```

* `superorphan` - matches "super" orphaned packages (dependencies no longer required by anything AND optional for nothing)
  ```
  qp w superorphan
  ```

  is equivalent to:
  ```
  qp where no:required-by and reason=dependency and no:optional-for
  ```

* `heavy` - matches packages 100MB and larger
  ```
  qp w superorphan
  ```

  is equivalent to:
  ```
  qp where size=100MB:
  ```

* `light` - matches packages 1MB and smaller
  ```
  qp w light
  ```

  is equivalent to:
  ```
  qp where size=1MB:
  ```

these macros can be combined with other queries as usual:

```
qp w orphan and size=100KB:
qp w not superorphan and not name=gtk
```

#### query examples

```
qp w name==bash and has:depends
qp where size=100MB:2GB
qp w name=python,cmake,yazi
qp w updated==2024-01-01
qp where q name=vim or name=nvim p and not has:conflicts
qp w not arch=x86_64
qp w q has:depends or has:required-by p and not reason=explicit
```

#### field types

| field type | description |
|------------|-------------|
| string | matches textual fields. used for fields like name, license, description, etc. <br> can take a comma-separated list |
| range | matches numerical or time-based fields across a range. <br> supports full ranges (start:end), open-ended ranges (start: / :end), or exact values |
| relation | matches fields that contain relationships to other packages (e.g., dependencies, conflicts, provides) <br> can take a comma-separated list |

#### fields and their types

| field name | field type |
|------------|------------|
| updated | range |
| built | range |
| size | range |
| name | string |
| reason | string |
| version | string |
| origin | string |
| arch | string |
| license | string |
| pkgbase | string |
| description | string |
| url | string
| validation | string |
| pkgtype | string |
| packager | string |
| groups | string |
| conflicts | relation |
| replaces | relation |
| depends | relation |
| optdepends | relation |
| required-by | relation |
| optional-for | relation |
| provides | relation |

### JSON output

`--output json` outputs the package data as structured JSON instead of a table. this can be useful for scripts or automation.

example:
```
qp select all where name=gtk3 --output json
```

`gtk3` is one of the few packages that actually has all the fields populated.

output format:
```json
[
  {
    "installTimestamp": 1743448253,
    "buildTimestamp": 1741400060,
    "size": 58266727,
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

### tips & tricks

- multiple short commands are supported using space separation (e.g. `s`, `w`, `l`, `o`), but **cannot** be combined as `swo` or `-swo`. use them like this:
  ```
  qp w name yay
  qp s name,size w name=vim o updated:asc l 10 # full query with shorthand
  ```

- group queries with `q ... p` to clarify order of operations:
  ```
  qp w q name=git or name=gh p and has:depends 
  ```

- the relations table columns can be lengthy. packages like `glibc` are required by thousands of packages. to improve readability, pipe the output to tools like `moar` or `less` (i prefer `moar`, but `less` is usually pre-installed):
  ```
  qp select name,depends | less
  qp s name,depends | moar
  ```

- options that take arguments can be used in the `--<option>=<value>` form:
  ```
  qp select name,updated --limit=100
  qp s name,updated o name
  ```

  boolean flags can be explicitly set using `--<option>=true` or `--<option>=false`:
  ```
  qp --no-headers=true --no-progress=true
  ```

  arguments to queries can be quoted if they contain special characters or spaces:
  ```
  qp where description="for tree-sitter"
  ```

- the `--no-headers` flag is useful when processing output in scripts. It removes the header row, making it easier to parse package lists with tools like `awk`, `sed`, or `cut`:
  ```
  qp --no-headers select name,size | awk '{print $1, $2}'
  ```

  **note**: `--no-progress` is automatically set to `true` in non-interactive environments, so you can pipe into programs like `cat`, `grep`, or `less` without issue.


### examples

 1. show the last 10 installed packages  
   ```
   qp limit 10
   ```

 2. show all explicitly installed packages  
   ```
   qp where reason=explicit limit all
   ```

 3. show only dependencies updated on a specific date  
   ```
   qp where reason=dependency and update=2025-03-01
   ```

 4. show all packages sorted alphabetically by name
   ```
   qp order name limit all
   ```

 5. search for packages that contain a GPL license 
   ```
   qp where license=gpl
   ```

 6. show packages installed between January 1, 2025, and January 5, 2025  
   ```
   qp where updated=2025-01-01:2025-01-05
   ```

 7. sort all packages by their license, displaying name and license  
   ```
   qp select name,license order license limit all
   ```

 8. show the 20 most recently installed packages larger than 20MB  
   ```
   qp where size=20MB: limit 20
   ```

 9. show packages between 100MB and 1GB installed up to February 27, 2025
   ```
   qp where size=100MB:1GB and updated=:2025-02-27
   ```

10. show all packages sorted by size in descending order, installed after January 1, 2025
   ```
   qp where updated=2025-01-01: order size:desc limit all
   ```

11. search for installed packages containing "python"  
   ```
   qp where name=python
   ```

12. search for explicitly installed packages containing "lib" that are between 10MB and 1GB in size
   ```
   qp where reason=explicit and name=lib and size=10MB:1GB
   ```

13. search for packages with names containing "linux" installed between January 1 and March 30, 2025
   ```
   qp where name=linux and updated=2025-01-01:2025-03-30
   ```

14. search for packages containing "gtk" installed after January 1, 2025, and at least 5MB in size
   ```
   qp where name=gtk and updated=2025-01-01: and size=5MB:
   ```

15. show packages with name, version, and size
   ```
   qp select name,version,size
   ```

16. show package names, descriptions, and dependencies with `less` for readability
   ```
   qp select name,depends,description | less
   ```

17. output package data in JSON format
   ```
   qp --output json
   ```

18. save all explicitly installed packages to a JSON file
   ```
   qp where reason=explicit --output json > explicit-packages.json
   ```

19. output all packages sorted by size (descending) in JSON
   ```
   qp order size:desc limit all --output json
   ```

20. output JSON with specific fields
   ```
   qp select name,version,size --output json
   ```

21. show all available package details for all packages
   ```
   qp select all limit all
   ```

22. output all packages with all fields in JSON format
   ```
   qp select all limit all --output json
   ```

23. show package names and sizes without headers for scripting
   ```
   qp select name,size --no-headers
   ```

24. show all packages required by `firefox`
   ```
   qp where required-by=firefox limit all
   ```

25. show all packages required by `gtk3` that are at least 50MB in size
   ```
   qp where required-by=gtk3 and size=50MB: limit all
   ```

26. show packages required by `vlc` and installed after January 1, 2025
   ```
   qp where required-by=vlc and updated=2025-01-01:
   ```

27. show all packages that have `glibc` as a dependency and are required by `ffmpeg`
   ```
   qp where depends=glibc and required-by=ffmpeg limit all
   ```

28. inclusively show packages that require `gcc` or `pacman`
   ```
   qp where required-by=gcc,pacman
   ```

29. show packages that provide `awk`
   ```
   qp where provides=awk
   ```

30. inclusively show packages that provide `rustc` or `python3`
   ```
   qp where provides=rustc,python3
   ```

31. show packages that conflict with `linuxqq`
   ```
   qp where conflicts=linuxqq
   ```

32. show packages that are built for the `aarch64` CPU architecture or any architecture
   ```
   qp where arch=aarch64,any
   ```

33. show all dependencies smaller than 500KB
   ```
   qp where reason=dependency and size=:500KB
   ```

34. show the 15 most recent explicitly installed packages
   ```
   qp where reason=explicit limit 15
   ```

35. show packages that contain "clang" in their description
   ```
   qp where description=clang
   ```

36. sort packages by their package base while showing their names and package bases, in reverse alphabetical order
   ```
   qp select name,pkgbase order pkgbase:desc
   ```

37. show packages that are exactly named "bash"
   ```
   qp where name==bash
   ```

38. show packages that have no dependencies
   ```
   qp where no:depends
   ```

39. show all packages installed via pipx
   ```
   qp where origin=pipx
   ```

## license

this project is licensed under GPL-3.0-only.

for use cases not compatible with the GPL, such as proprietary redistribution or integration/ingestion into ML/LLM systems, a separate commercial license is available. see LICENSE.commercial for details.
