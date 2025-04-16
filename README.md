# qp - query packages

`qp` is a CLI util, written in **Go** / **Golang**, for [arch linux](https://archlinux.org) and arch-based linux distros to query installed packages.

you can find installation instructions [here](#installation).

`qp` supports querying and sorting for install date, package name, install reason (explicit/dependency), size on disk, reverse dependencies (required by), dependencies, validation, description, replacements, conflicts, provisions, build date, packager, optional dependencies, reverse optional dependencies (optional for), package type and more. check [usage](#usage) for all available options.

![qp logo | query packages logo](https://gistcdn.githack.com/Zweih/9009d5c74eab8a5515a8a64a0495df32/raw/ef8a8ac3655fd3dee24494a3403867919d806b63/qp-logo_clean.svg)

[![AUR version - qp](https://img.shields.io/aur/version/qp?style=flat-square&logo=arch-linux&logoColor=1793d1&label=qp&color=1793d1)](https://aur.archlinux.org/packages/qp)
[![AUR version - qp-bin](https://img.shields.io/aur/version/qp-bin?style=flat-square&logo=arch-linux&logoColor=1793d1&label=qp-bin&color=1793d1)](https://aur.archlinux.org/packages/qp-bin)
[![AUR version - qp-git](https://img.shields.io/aur/version/qp-git?style=flat-square&logo=arch-linux&logoColor=1793d1&label=qp-git&color=1793d1)](https://aur.archlinux.org/packages/qp-git)

![GitHub Downloads](https://img.shields.io/github/downloads/Zweih/qp/total?style=for-the-badge&logo=github&label=Downloads%20since%202%2F4%2F2025&color=1793d1)

![Alt](https://repobeats.axiom.co/api/embed/a13406d103a649d70641774ee85e7a9983ccf96b.svg "Repobeats analytics image")

<details open>
<summary><strong>download and clone statistics</strong></summary>
<br>
 
graphs are generated daily with my other project, [repulse analytics](https://github.com/Zweih/repulse-analytics)

<img src="https://raw.githubusercontent.com/Zweih/repulse-analytics/refs/heads/repulse-traffic-graphs/total_downloads.png" alt="total downloads" width="400"/> <img src="https://raw.githubusercontent.com/Zweih/repulse-analytics/refs/heads/repulse-traffic-graphs/total_clones.png" alt="total clones" width="400"/>

</details>

this package is compatible with the following distributions:
 - [arch linux](https://archlinux.org)
 - [manjaro](https://manjaro.org/)
 - [steamOS](https://store.steampowered.com/steamos)
 - [garuda linux](https://garudalinux.org/)
 - [endeavourOS](https://endeavouros.com/)
 - [artix linux](https://artixlinux.org/)
 - the 50 other arch-based distros, as long as it has pacman installed 

## features

- list installed packages with install date/timestamps, dependencies, provisions, reverse dependencies (required by), size on disk, conflicts, replacements, architecture, license, description, build date, package base, package type, validation, packager, optional dependencies, reverse optional dependencies (optional for), groups, and version
- query by field existence
- query by explicitly installed packages
- query by packages installed as dependencies
- query by packages required by specified packages
- query by packages that depend upon specified packages
- query by packages that provide specified packages
- query by packages that conflict with specific packages
- query by packages that contain specific licenses
- query by a specific installation date or date range
- query by packages built with specified architectures
- query by package size or size range
- query by package names
- query by package description
- sort by installation date, package name, license, or by size on disk
- output as a table or JSON

## is it good?
[yes.](https://news.ycombinator.com/item?id=3067434)

## roadmap

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
| – | list possibly or confirmed stale/abandoned packages | – | self-referencing field |
| – | name exclusion query | – | streaming pipeline |
| – | short-args for queries | – | key/value output |
| – | XML output | – | package description sort |
| – | package base query | – | required-by sort |
| – | required-by count sort | – | dependency count sort |
| ✓ | build-date field | - | build-date filter |
| - | build-date sort | ✓ | pkgtype field |
| - | pkgtype filter | - | pkgtype sort |
| ✓ | architecture query | ✓ | groups field |
| ✓	| conflicts query | - | package description sort |
| ✓	| regenerate cache option | - | groups filter |
| ✓ | packager field | ✓ | optional dependency field |
| ✓ | sort by size on disk | - | conflicts sort |
| ✓ | validation field | - | validation sort |
| ✓ | reverse optional dependencies field (optional for) | - | optdepends installation indicator |
| - | optional-for query | - | separate field for optdepends reason |
| ✓ | fuzzy/strict querying | - | exclusion querying |
| ✓ | existence querying | - | depth querying |

## installation

### from AUR (**recommended**)
install using an [AUR helper](https://wiki.archlinux.org/title/AUR_helpers) like `yay`, `paru`, `aura`, etc.:
```bash
yay -Sy qp
```

if you prefer to install a pre-compiled binary* using the AUR, use the `qp-bin` package instead.

***note**: binaries are automatically, securely, and transparently compiled with github CI when a version release is created. you can audit the binary creation by checking the relevant github action for each release version.

for the latest (unstable) version from git w/ the AUR, use `qp-git`*.  

***note**: this is not recommended for most users

the cache is located under `/query-packages` at `$HOME/.cache/` or wherever you have `$XDG_HOME_CACHE` set to.

### building from source + manual installation
**note**: this packages is specific to arch-based linux distributions

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

## usage

```bash
qp [options]
```

### options
- `-l <number>` | `--limit <number>`: limit the amount of recent packages to display (default: 20)
- `-a` | `all`: show all installed packages (ignores `-l`)
- `-w <query>` | `--where <query>`: apply one or more filters to refine package results.
  - supported query types:
    - **string match** -> `field=value` (fuzzy) or `field==value` (strict)
    - **range match** -> `field=start:end` (fuzzy) or `field==start:end` (strict)
        - supports full ranges (`start:end`), open-ended ranges (`start:` or `:end`), and exact values (`value`)
        - works with `date` and `size`
    - **existence check** -> `has:field` or `not:field`
  - this flag can be used multiple times and mixed freely
  - [see all available query fields below](#available-queries)
- `-O <field>:<direction>` | `--order <field>:<direction>`: sort results ascending or descending (default sort is `date:asc`):
  - `date`    -> sort by installation date
  - `name`    -> sort alphabetically by package name
  - `size`    -> sort by package size on disk
  - `license` -> sort alphabetically by package license
  - `pkgbase` -> sort alphabetically by base package
- `--no-headers`: omit column headers in table output (useful for scripting)
- `-s <list>` | `--select <list>`: comma-separated list of fields to display (cannot use with `--select-all` or `--select-add`)
- `-S <list>` | `--select-add <list>`: comma-separated list of fields to add to defaults or `--select-all`
- `-A` | `--select-all`: output all available fields (overrides defaults)
- `--full-timestamp`: display the full timestamp (date and time) of package install/build instead of just the date
- `--json`: output results in JSON format (overrides table output and `--full-timestamp`)
- `--no-progress`: force no progress bar outside of non-interactive environments
- `--no-cache`: disable cache loading/saving and force fresh package data loading
- `--regen-cache`: disable cache loading, force fresh package data loading, and save fresh cache
- `-h` | `--help`: print help info

### querying with `--where`
the `--where` (short: `-w`) flag is the core of qp's flexible query system. you can use it multiple times per command to combine different filters.

#### query types
all queries that take words as arguments can also take a comma-separated list.

each `--where` query supports one of the following:

- **string match**
  - `field=value` -> fuzzy match
  - `field==value` -> strict match
  - applies to string-based fields like `name`, `license`, `description`, etc.
- **range match**
  - `field=start:end` -> fuzzy match
  - `field==start:end` -> strict match
  - works with `date` and `size`
  - supports:
    - full ranges: `start:end`
    - open-ended ranges: `start:` or `:end`
    - exact values: `value`
- **existence check**
  - `has:field` — field must exist or be non-empty
  - `not:field` — field must be missing or empty

#### match types

|field type | fuzzy | strict | 
|-----------|-------|--------|
| strings & relations | substring match (case-insensitive) | exact character match (case-insensitive) |
| dates | matches by day (ignores time) | matches exact timestamp (to the second) |
| size  | ±0.3% byte tolerance (approximate) | matches exact byte size |

For example:
  - `name=gtk` matches `gtk3`, `libgtk`, etc. (fuzzy)
  - `name==gtk` only matches a package named exactly `gtk`

#### query examples
```bash
qp -w size=100MB:1GB          # size range (fuzzy)
qp -w date==2024-01-01        # exact install date
qp -w name=firefox            # fuzzy name match
qp -w name==bash              # strict name match
qp -w reason=explicit         # packages installed explicitly
qp -w has:depends             # must have dependencies
qp -w not:conflicts           # must not conflict with anything
```

#### available queries
| query type  | syntax | description |
|-------------|--------|-------------|
| **date** | `date=<value>` | query by installation date. supports exact dates, ranges (`YYYY-MM-DD:YYYY-MM-DD`), and open-ended ranges (`YYYY-MM-DD:` or `:YYYY-MM-DD`) |
| **size** | `size=<value>` | query by package size on disk. supports exact values (`10MB`), ranges (`10MB:1GB`), and open-ended ranges (`:500KB`, `1GB:`) |
| **name** | `name=<package>` / <br> `name=<package-1>,<package-2>,<etc>` | query by package name |
| **installation reason** | `reason=explicit` / `reason=dependencies` | query packages by installation reason: explicitly installed or installed as a dependency |
| **architecture** | `arch=<architecture>` / <br> `arch=<architecture-1>,<architecture-2>,<etc>` | query by packages that are built for the specified architectures <br> **note**: "any" is a separate architecture category |
| **license** | `license=<license>` / <br> `license=<license-1>,<license-2>,<etc>` | query by license name (substring match) |
| **description** | `description=<string>` / <br> `description=<string-1>,<string-2>,<etc>` | query by package description |
| **conflicts** | `conflicts=<package>` / <br> `conflicts=<package-1,<package-2>,<etc>` | query by packages that conflict with the specified packages |
| **depends** | `depends=<package>` / <br> `depends=<package-1>,<package-2>,<etc>` | query by packages that have the specified packages as dependencies |
| **required by** | `required-by=<package>` / <br> `required-by=<package-1>,<package-2>,<etc>` | query by packages that are required by the specified packages |
| **provides** | `provides=<package>` / <br> `provides=<package-1>,<package-2>,<etc>` | query by package that provide the specified packages/libraries |

### available fields for selection
- `date` - installation date of the package
- `build-date` - date the package was built
- `size` - package size on disk
- `name` - package name
- `reason` - installation reason (explicit/dependency)
- `version` - installed package version
- `arch` - architecture the package was built for (e.g., x86_64, aarch64, any)
- `license` - package software license
- `pkgbase` - name of the base package used to group split packages; for non-split packages, it is the same as the package name. 
- `description` - package description
- `url` - the URL of the official site of the software being packaged
- `validation` - package integrity validation method (e.g., sha256", "pgp")
- `packager` - person/entity who built the package (if available)
- `pkgtype` - package type (pkg, split, debug, source, unknown*)
    - ***note**: older packages may show "unknown" pkgtype if built before pacman introduced XDATA
- `groups` - package groups or categories (e.g., base, gnome, xfce4)
- `conflicts` - list of packages that conflict, or cause problems, with the package
- `replaces` - list of packages that are replaced by the package
- `depends` - list of dependencies
- `optdepends` - list of optional dependencies
- `required-by` - list of packages required by the package and are dependent
- `optional-for` - list of packages that optionally depend on the package (optionally dependent)
- `provides` - list of alternative package names or shared libraries provided by package

### JSON output
the `--json` flag outputs the package data as structured JSON instead of a table. this can be useful for scripts or automation.

example:
```bash
qp -Aw name=gtk3 --json
```

`gtk3` is one of the few packages that actually has all the fields populated.

output format:
```json
[
  { 
    "installTimestamp": 1743448253,
    "buildTimestamp": 1741400060,
    "size": 58266727,
    "pkgtype": "split",
    "name": "gtk3",
    "reason": "dependency",
    "version": "1:3.24.49-1",
    "arch": "aarch64",
    "license": "LGPL-2.1-or-later",
    "pkgbase": "gtk3",
    "description": "GObject-based multi-platform GUI toolkit",
    "url": "https://www.gtk.org/",
    "validation": "pgp",
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

- when using multiple short flags at once (e.g. `-aw` or `-Al`), the flags like `-w`, `-l`, and `-s` must be last as they consume the next argument.
this follows standard unix-style flag parsing, where positional arguments (like numbers and strings)
are treated as separate parameters.
  
  invalid:
  ```bash
  qp -wa name=yay  # incorrect usage 
  ```
  valid:
  ```bash
  qp -aw name=yay  # correct usage
  ```

- the `depends`, `provides`, and `required-by` table columns can be lengthy. packages like `glibc` are required by thousands of packages. to improve readability, pipe the output to tools like `moar` or `less` (i prefer `moar`, but `less` is usually pre-installed):
  ```bash
  qp -s name,depends | less
  ```
- all options that take an argument can also be used in the `--<flag>=<argument>` format:
  ```bash
  qp --select-add=name --limit=100
  qp -s=date,name,version -O=name
  ```
  boolean flags can also be explicitly set using `--<flag>=true` or `--<flag>=false`:
  ```bash
  qp --no-headers=true --no-progress=true
  ```
  string arguments can also be surrounded with quotes or double-quotes:
  ```bash
  qp --order="name" -w name="vim"
  ```

  this can be useful for scripts and automation where you might want to avoid any and all ambiguity.

  **note**: `--no-progress` is automatically set to `true` when in a non-interactive environment, so you can pipe `|` into programs like `cat`, `grep`, or `less` without issue

- the `--no-headers` flag is useful when processing output in scripts. It removes the header row, making it easier to parse package lists with tools like `awk`, `sed`, or `cut`:
  ```bash
  qp --no-headers --select name,size | awk '{print $1, $2}'
  ```

### examples

 1. show the last 10 installed packages 
   ```bash
   qp -l 10
   ```
 2. show all explicitly installed packages
   ```bash
   qp -aw reason=explicit
   ```
 3. show only dependencies installed on a specific date
   ```bash
   qp -w reason=dependency -w date=2025-03-01
   ```
 4. show all packages sorted alphabetically by name
   ```bash
   qp -aO name
   ```
 5. search for packages that contain a GPL license
   ```bash
   qp -w license=gpl
   ```
 6. show packages installed between january 1, 2025, and january 5, 2025
   ```bash
   qp -w date=2025-01-01:2025-01-05
   ```
 7. sort all packages by their license, displaying name and license
   ```bash
   qp -aO license -s name,license
   ```
 8. show the 20 most recently installed packages larger than 20MB
   ```bash
   qp -w size=20MB: -l 20
   ```
 9. show packages between 100MB and 1GB installed up to february 27, 2025
   ```bash
   qp -w size=100MB:1GB -w date=:2025-02-27
   ```
10. show all packages sorted by size in descending order, installed after january 1, 2025
   ```bash
   qp -a --order size:desc -w date=2025-01-01:
   ```
11. search for installed packages containing "python
   ```bash
   qp -w name=python
   ```
12. search for explicitly installed packages containing "lib" that are between 10MB and 1GB in size
   ```bash
   qp -w reason=explicit -w name=lib -w size=10MB:1GB
   ```
13. search for packages with names containing "linux" installed between january 1 and march 30, 2025
   ```bash
   qp -w name=linux -w date=2025-01-01:2025-03-30
   ```
14. search for packages containing "gtk" installed after january 1, 2025, and at least 5MB in size
   ```bash
   qp -w name=gtk -w date=2025-01-01: -w size=5MB:
   ```
15. show packages with name, version, and size
   ```bash
   qp -s name,version,size
   ```
16. show package names, descriptions, and dependencies with `less` for readability
   ```bash
   qp --select name,depends,description | less
   ```
17. output package data in JSON format
   ```bash
   qp --json
   ```
18. save all explicitly installed packages to a JSON file
   ```bash
   qp -w reason=explicit --json > explicit-packages.json
   ```
19. output all packages sorted by size (descending) in JSON
   ```bash
   qp --json -a -O size:desc
   ```
20. output JSON with specific fields
   ```bash
   qp --json -s name,version,size
   ```
21. show all available package details for all packages
   ```bash
   qp -aA
   ```
22. output all packages with all fields in JSON format
   ```bash
   qp -aA --json
   ```
23. show package names and sizes without headers for scripting
   ```bash
   qp --no-headers -s name,size
   ```
24. show all packages required by `firefox`
   ```bash
   qp -a -w required-by=firefox
   ```
25. show all packages required by `gtk3` that are at least 50MB in size
   ```bash
   qp -a -w required-by=gtk3 -w size=50MB:
   ```
26. show packages required by `vlc` and installed after january 1, 2025 
   ```bash
   qp -w required-by=vlc -w date=2025-01-01:
   ```
27. show all packages that have `glibc` as a dependency and are required by `ffmpeg`
   ```bash
   qp -a -w depends=glibc -w required-by=ffmpeg
   ```
28. inclusively show packages that require `gcc` or `pacman`:
   ```bash
   qp -w required-by=base-devel,gcc
   ```
29. show packages that provide `awk`:
   ```bash
   qp -w provides=awk
   ```
30. inclusively show packages that provide `rustc` or `python3`:
   ```bash
   qp -w provides=rustc,python3
   ```
31. show packages that conflict with `linuxqq`:
   ```bash
   qp -w conflicts=linuxqq
   ```
32. show packages that are built for the `aarch64` CPU architecture or any architecture (non-CPU-specific):
   ```bash
   qp -w arch=aarch64,any
   ```
33. show all dependencies smaller than 500KB  
   ```bash
   qp -w reason=dependencies -w size=:500KB
   ```
34. show the 15 most recent explicitly installed packages
   ```bash
   qp -w reason=explicit -l 15
   ```
35. show packages that contain "clang" in their description:
   ```bash
   qp -w description=clang
   ```
36. sort packages by their package base while showing their names and package bases, in reverse alphabetical order:
   ```bash
   qp -O pkgbase:desc -s name,pkgbase
   ```
37. show packages that are exactly named "bash":
   ```bash
   qp -w name==bash
   ```

## license
this project is licensed under GPL-3.0-only.

for use cases not compatible with the GPL, such as proprietary redistribution or integration/ingestion into ML/LLM systems, a separate commercial license is available. see LICENSE.commercial for details.
 
