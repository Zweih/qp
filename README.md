# qp - query packages

`qp` is a CLI util, written in **Go** / **Golang**, for [arch linux](https://archlinux.org) and arch-based linux distros to query installed packages.

you can find installation instructions [here](#installation).

`qp` supports querying/sorting for install date, package name, install reason (explicit/dependency), size on disk, reverse dependencies, dependency requirements, description, and more. check [usage](#usage) for all available options.

![AUR status - qp](https://img.shields.io/aur/version/qp?style=flat-square&label=AUR%20-%20qp&link=https%3A%2F%2Faur.archlinux.org%2Fpackages%2Fqp)
![AUR status - qp-bin](https://img.shields.io/aur/version/qp-bin?style=flat-square&label=AUR%20-%20qp-bin&link=https%3A%2F%2Faur.archlinux.org%2Fpackages%2Fqp-bin)
![AUR status - qp-git](https://img.shields.io/aur/version/qp-git?style=flat-square&label=AUR%20-%20qp-git&link=https%3A%2F%2Faur.archlinux.org%2Fpackages%2Fqp-git)

![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/Zweih/qp/total?style=for-the-badge&logo=archlinux&label=Downloads%20Since%202%2F4%2F2025&color=%23007ec6)

![Alt](https://repobeats.axiom.co/api/embed/7a20b73b689d45d678001c582a9d1f124dca31ba.svg "Repobeats analytics image")

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

- list installed packages with date/timestamps, dependencies, provisions, requirements, size on disk, conflicts, architecture, license, description, and version
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

## why is it called qp if it works with other AUR helpers?
because yay is my preferred AUR helper and the name has a good flow.

## is it good?
[yes.](https://news.ycombinator.com/item?id=3067434)

## roadmap

- [x] rewrite in golang
- [x] additional queries
- [ ] list possibly or confirmed stale/abandoned packages
- [x] sort by size on disk
- [x] protobuf caching (127% speed boost)
- [ ] dependency graph
- [x] concurrent querying
- [x] query by size on disk
- [x] asynchronous progress bar
- [x] channel-based aggregation
- [x] concurrent sorting
- [x] query by package name
- [x] package version field
- [x] query by date range
- [x] concurrent file reading (200% speed boost)
- [x] remove expac as a dependency (300% speed boost)
- [x] package provisions
- [x] optional full timestamp 
- [x] add CI to release binaries
- [x] remove go as a dependency
- [x] query by range of size on disk
- [x] user defined columns
- [x] dependencies of each package (dependency field)
- [x] reverse-dependencies of each package (required-by field)
- [x] package description field
- [x] package URL field
- [x] package architecture field
- [x] package conflicts field
- [x] conflicts query
- [ ] name exclusion query
- [ ] self-referencing field
- [x] JSON output
- [x] no-headers option
- [x] provides query
- [x] depends query
- [x] all-columns option
- [x] required-by query
- [ ] key/value output
- [x] list of packages for package queries
- [x] config dependency injection for testing
- [ ] required-by count sort
- [x] optimize file reading (28% speed boost)
- [x] metaflag for all queries
- [x] package license field
- [ ] XML output
- [x] use chunked channel-based concurrent querying (12% speed boost) 
- [ ] short-args for queries
- [x] license sort
- [ ] packager field
- [ ] streaming pipeline
- [x] architecture query
- [x] optimize query order (4% speed boost)
- [ ] dependency count sort
- [x] license query
- [ ] optional dependency field
- [x] improve sorting efficiency (8% speed boost)
- [ ] package base field
- [x] package description query

## installation

### from AUR (**recommended**)
install using [AUR helper](https://wiki.archlinux.org/title/AUR_helpers) like `yay`, `paru`, `aura`, etc.:
```bash
yay -Sy qp
```

if you prefer to install a pre-compiled binary* using the AUR, use the `qp-bin` package instead.

***note**: binaries are automatically, securely, and transparently compiled with github CI when a version release is created. you can audit the binary creation by checking the relevant github action for each release version.

for the latest (unstable) version from git w/ the AUR, use `qp-git`*.  

***note**: this is not recommended for most users 


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
- `-w <field>=<value>` | `--where <field>=<value>`: apply multiple queries for a flexible query system. can be used multiple times. example:
  - `--where size=10MB:1GB`    -> query by size range
  - `--where date=2024-01-01:` -> query by installation date
  - `--where reason=explicit`  -> query by explicit installations
  - `--where name=firefox`     -> query by package names that contain `firefox`
  - `--where required-by=vlc`  -> query by packages required by `vlc`
  - `--where depends=glibc`    -> query by packages that depend on `glibc`
  - `--where conflicts=sdl2`   -> query by packages that conflict with `sdl2`
  - `--where arch=x86_64`      -> query by packages that are built for `x86_64` CPUs
  - `--where license=GPL`      -> query by package licenses that contain `GPL`
  - `--where description="linux kernel"` -> query by package descriptions that contain "linux kernel" 
- `-O <field>:<direction>` | `--order <field>:<direction>`: sort results ascending or descending (default sort is `date:asc`):
  - `date`    -> sort by installation date
  - `name`    -> sort alphabetically by package name
  - `size`    -> sort by package size on disk
  - `license` -> sort alphabetically by package license
- `--no-headers`: omit column headers in table output (useful for scripting)
- `-s <list>` | `--select <list>`: comma-separated list of fields to display (cannot use with `--select-all` or `--select-add`)
- `-S <list>` | `--select-add <list>`: comma-separated list of fields to add to defaults or `--select-all`
- `-A` | `--select-all`: output all available fields (overrides defaults)
- `--full-timestamp`: display the full timestamp (date and time) of package installations instead of just the date
- `--json`: output results in JSON format (overrides table output and `--full-timestamp`)
- `--no-progress`: force no progress bar outside of non-interactive environments
- `-h` | `--help`: print help info

### querying with `--where`
the `--where` (short-flag: `-w`) flag allows for powerful querying of installed packages. queries can be combined by using multiple query flags. 

all queries that take package/library/program names as arguments can also take a comma-separated list. this applies to the queries `name`, `depends`, and `required-by`. 

short-flag queries and long-flag queries can be combined.

#### available queries
| query type  | syntax | description |
|-------------|--------|-------------|
| **license** | `license=<license>` / <br> `license=<license-1>,<license-2>,<etc` | query by license name (substring match) |
| **date** | `date=<value>` | query by installation date. supports exact dates, ranges (`YYYY-MM-DD:YYYY-MM-DD`), and open-ended ranges (`YYYY-MM-DD:` or `:YYYY-MM-DD`) |
| **required by** | `required-by=<package>` / <br> `required-by=<package-1>,<package-2>,<etc>` | query by packages that are required by the specified packages |
| **depends** | `depends=<package>` / <br> `depends=<package-1>,<package-2>,<etc>` | query by packages that have the specified packages as dependencies |
| **provides** | `provides=<package>` / <br> `provides=<package-1>,<package-2>,<etc>` | query by package that provide the specified packages/libraries |
| **conflicts** | `conflicts=<package>` / <br> `conflicts=<package-1,<package-2>,<etc>` | query by packages that conflict with the specified packages |
| **architecture** | `arch=<architecture>` / <br> `arch=<architecture-1>,<architecture-2>,<etc>` | query by packages that are built for the specified architectures <br> **note**: "any" is a separate architecture category |
| **name** | `name=<package>` / <br> `name=<package-1>,<package-2>,<etc>` | query by package name (substring match) |
| **installation reason** | `reason=explicit` / `reason=dependencies` | query packages by installation reason: explicitly installed or installed as a dependency |
| **size** | `size=<value>` | query by package size on disk. supports exact values (`10MB`), ranges (`10MB:1GB`), and open-ended ranges (`:500KB`, `1GB:`) |
| **description** | `description=<string>` / <br> `description=<string-1>,<string-2>,<etc>` | query by package description (substring match) |

### available fields
- `date` - installation date of the package
- `name` - package name
- `reason` - installation reason (explicit/dependency)
- `size` - package size on disk
- `version` - installed package version
- `depends` - list of dependencies (output can be long)
- `required-by` - list of packages required by the package and are dependent (output can be long) 
- `provides` - list of alternative package names or shared libraries provided by package (output can be long)
- `conflicts` - list of packages that conflict, or cause problems, with the package
- `arch` - architecture the package was built for (e.g., x86_64, aarch64, any)
- `license` - package software license
- `url` - the URL of the official site of the software being packaged
- `description` - package description

### JSON output
the `--json` flag outputs the package data as structured JSON instead of a table. this can be useful for scripts or automation.

example:
```bash
qp -Aw name=tinysparql --json
```

`tinysparql` is one of the few packages that actually has all the fields populated.

output format:
```json
[
  {
    "timestamp": 1739294218,
    "size": 4385422,
    "name": "tinysparql",
    "reason": "dependency",
    "version": "3.8.2-2",
    "arch": "aarch64",
    "license": "GPL-2.0-or-later",
    "url": "https://tinysparql.org/",
    "description": "Low-footprint RDF triple store with SPARQL 1.1 interface",
    "depends": [
      "avahi",
      "gcc-libs",
      "glib2",
      "glibc",
      "icu",
      "json-glib",
      "libsoup3",
      "libstemmer",
      "libxml2",
      "sqlite"
    ],
    "provides": [
      "tracker3=3.8.2",
      "libtinysparql-3.0.so=0-64"
    ],
    "conflicts": [
      "tracker3<=3.7.3-2"
    ]
  }
]
```

### tips & tricks

- when using multiple short flags, the -l flag must be last since it consumes the next argument.
this follows standard unix-style flag parsing, where positional arguments (like numbers)
are treated as separate parameters.
  
  invalid:
  ```bash
  qp -wa name=yay  # incorrect usage 
  ```
  valid:
  ```bash
  qp -aw name=yay  # correct usage
  ```

- the `depends`, `provides`, `required-by` columns output can be lengthy, packages like `glibc` are required by thousands of packages. to improve readability, pipe the output to tools like `less` or `moar` (i prefer `moar`, but `less` is a core-util):
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
  qp --no-headers --columns name,size | awk '{print $1, $2}'
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
22. output all packages with all columns/fields in JSON format
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

## license
this project is licensed under GPL-3.0-only.

for use cases not compatible with the GPL, such as proprietary redistribution or integration/ingestion into ML/LLM systems, a separate commercial license is available. see LICENSE.commercial for details.
 
