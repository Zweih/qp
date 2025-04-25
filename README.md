# qp - query packages

`qp` is a CLI util, written in **Go** / **Golang**, for [arch linux](https://archlinux.org) and arch-based linux distros to query installed packages.

you can find installation instructions [here](#installation).

`qp` supports querying with full boolean logic for package metadata, dependency relations, and more.

check [features](#features) to find out more.

check [usage](#usage) for all available commands + options.

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
 - [mabox linux](https://maboxlinux.org/)
 - [artix linux](https://artixlinux.org/)
 - the 50 other arch-based distros, as long as it has pacman installed

non-arch distros are planned!

## features

- list installed packages by numerous fields
  - see all available fields for selection [here](#available-fields-for-selection)
- query by:
  - field existence
  - explicitly installed packages
  - packages installed as dependencies, required by specified packages, depend upon specified packages, provide specified packages, conflict with specific packages
  - packages that contain specific licenses
  - installation/build date or date range
  - packages built with specified architectures
  - package size or size range
  - package names, descriptions, and packagers
  - learn more about querying [here](#querying-with-where)
- sort by: 
  - installation date and build date
  - package name
  - license
  - size on disk
  - package base and package type
- output as:
  - table
  - JSON


learn about usage [here](#usage)

learn about installation [here](#installation)

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
| ✓  | groups query | – | streaming pipeline |
| – | short-args for queries | – | key/value output |
| – | XML output | – | package description sort |
| ✓ | package base query | – | required-by sort |
| – | optdepends sort | – | depends sort |
| ✓ | build-date field | ✓ | build-date query |
| ✓ | build-date sort | ✓ | pkgtype field |
| ✓ | url query | ✓ | pkgtype sort |
| ✓ | architecture query | ✓ | groups field |
| ✓	| conflicts query | - | package description sort |
| ✓	| regenerate cache option | ✓ | validation query |
| - | url sort | - | groups sort |
| ✓ | packager field | ✓ | optional dependency field |
| ✓ | sort by size on disk | - | conflicts sort |
| - | optional-for sort | - | provides sort |
| ✓ | validation field | - | validation sort |
| - | packager sort | - | architecture sort |
| - | reason sort | - |version sort |
| ✓ | reverse optional dependencies field (optional for) | - | optdepends installation indicator |
| ✓ | optional-for query | - | separate field for optdepends reason |
| ✓ | fuzzy/strict querying | ✓ | existence querying |
| ✓ | existence querying | - | depth querying |
| ✓ | pkgtype query | ✓ | optdepends query |
| ✓ | packager query | - | origin field |
| - | origin sort | - | origin query |
| ✓ | command-based syntax | ✓ | full boolean logic |
| ✓ | abstract syntax tree | ✓ | directed acyclical graph for filtering |
| - | user-defined macros | ✓ | parentetical (grouping) logic |
| ✓ | limit from end | ✓ | limit from middle |
| - | replaces sort | - | built-in macros |
| - | query explaination | - | user configuration file |

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
qp [command] [args] [options]
```

### commands

- `select <list>` | `s <list>`: comma-separated list of fields to display
  - `select all | s all` will act as a list of all available fields
  - `select default | s default` will act as a list of all default fields
  - use `select default,version` to list default fields + version
  - use `select all,version` to list default fields + version
  - [see fields available for selection](#available-selectors)
- `where <query>` | `w <query>`: apply one or more queries to refine package results.
  - supported query types:
    - **string match** -> `field=value` (fuzzy) or `field==value` (strict)
    - **range match** -> `field=start:end` (fuzzy) or `field==start:end` (strict)
        - supports full ranges (`start:end`), open-ended ranges (`start:` or `:end`), and exact values (`value`)
        - works with `date` and `size`
    - **existence check** -> `has:field` or `no:field`
  - use `and`, `or`, `not`, `q ... p` to build complex filters
  - learn more about querying [here](#querying-with-where)
- `order <field>:<direction>` | `o <field>:<direction>`: sort results ascending or descending
  - default sort is `date:asc`:
  - [see fields avaialble for sorting](#available-sorts)
- `limit <number>` | `l <number>`: limit the amount of packages to display (default: 20)
  - `limit all | l all`: display all packages
  - `limit end:<number>`: display last `n` packages
  - `limit mid:<number>`: display middle `n` packages

### options

- `--no-headers`: omit column headers in table output (useful for scripting)
- `--full-timestamp`: display the full timestamp (date and time) of package install/build instead of just the date
- `--json`: output results in JSON format (overrides table output and `--full-timestamp`)
- `--no-progress`: force no progress bar outside of non-interactive environments
- `--no-cache`: disable cache loading/saving and force fresh package data loading
- `--regen-cache`: disable cache loading, force fresh package data loading, and save fresh cache
- `-h` | `--help`: print help info

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
  - works with `date` and `size`
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

#### query examples

```
qp w name==bash and has:depends
qp where size=100MB:2GB
qp w name=python,cmake,yazi
qp w date==2024-01-01
qp where q name=vim or name=nvim p and not has:conflicts
qp w not arch=x86_64
qp w q has:depends or has:required-by p and not reason=explicit
```

#### query types

| field type | description |
|------------|-------------|
| string | matches textual fields. used for fields like name, license, description, etc. <br> can take a comma-separated list |
| range | matches numerical or time-based fields across a range. <br> supports full ranges (start:end), open-ended ranges (start: / :end), or exact values |
| relation | matches fields that contain relationships to other packages (e.g., dependencies, conflicts, provides) <br> can take a comma-separated list |

#### available queries

| field name | field type |
|------------|------------|
| date | range |
| build-date | range |
| size | range |
| name | string |
| reason | string |
| version | string |
| pkgtype | string |
| arch | string |
| license | string |
| pkgbase | string |
| description | string |
| url | string
| validation | string |
| packager | string |
| groups | string |
| conflicts | relation |
| replaces | relation |
| depends | relation |
| optdepends | relation |
| required-by | relation |
| optional-for | relation |
| provides | relation |

### available selectors

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
- `validation` - package integrity validation method (e.g., sha256, pgp)
- `pkgtype` - package type (pkg, split, debug, src)
    - ***note**: older packages may have no pkgtype if built before pacman introduced XDATA
- `packager` - person/entity who built the package (if available)
- `groups` - package groups or categories (e.g., base, gnome, xfce4)
- `conflicts` - list of packages that conflict, or cause problems, with the package
- `replaces` - list of packages that are replaced by the package
- `depends` - list of dependencies
- `optdepends` - list of optional dependencies
- `required-by` - list of packages required by the package and are dependent
- `optional-for` - list of packages that optionally depend on the package (optionally dependent)
- `provides` - list of alternative package names or shared libraries provided by package

### available sorts

- `date`
- `build-date`
- `name`
- `license`
- `size`
- `pkgtype`
- `pkgbase`

### JSON output

the `--json` flag outputs the package data as structured JSON instead of a table. this can be useful for scripts or automation.

example:
```
qp select all where name=gtk3 --json
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
    "pkgtype": "split",
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

- multiple short commands are supported using space separation (e.g. `s`, `w`, `l`, `o`), but **cannot** be combined as `swo` or `-swo`. use them like this:
  ```
  qp w name yay
  qp s name,size w name=vim o date:asc l 10 # full query with shorthand
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
  qp select name,date --limit=100
  qp s name,date o name
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

 3. show only dependencies installed on a specific date  
   ```
   qp where reason=dependency and date=2025-03-01
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
   qp where date=2025-01-01:2025-01-05
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
   qp where size=100MB:1GB and date=:2025-02-27
   ```

10. show all packages sorted by size in descending order, installed after January 1, 2025
   ```
   qp where date=2025-01-01: order size:desc limit all
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
   qp where name=linux and date=2025-01-01:2025-03-30
   ```

14. search for packages containing "gtk" installed after January 1, 2025, and at least 5MB in size
   ```
   qp where name=gtk and date=2025-01-01: and size=5MB:
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
   qp --json
   ```

18. save all explicitly installed packages to a JSON file
   ```
   qp where reason=explicit --json > explicit-packages.json
   ```

19. output all packages sorted by size (descending) in JSON
   ```
   qp order size:desc limit all --json
   ```

20. output JSON with specific fields
   ```
   qp select name,version,size --json
   ```

21. show all available package details for all packages
   ```
   qp select all limit all
   ```

22. output all packages with all fields in JSON format
   ```
   qp select all limit all --json
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
   qp where required-by=vlc and date=2025-01-01:
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

## license

this project is licensed under GPL-3.0-only.

for use cases not compatible with the GPL, such as proprietary redistribution or integration/ingestion into ML/LLM systems, a separate commercial license is available. see LICENSE.commercial for details.
