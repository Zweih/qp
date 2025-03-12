# yaylog

`yaylog` is a CLI util, written in **Go** / **Golang**, for [arch linux](https://archlinux.org) and arch-based linux distros to sort/filter installed packages.

despite the name, it's not limited to `yay` and works with any package manager that uses ALPM; so it can be used with `pacman`, `yay`, `paru`, `aura`, `pamac`, and even `yaourt` if you're somehow still using it.

`yaylog` supports optional filters/sorting for install date, package name, install reason (explicit/dependency), size on disk, reverse dependencies, dependency requirements, and more. check [usage](#usage) for all available options.

[![Packaging status](https://repology.org/badge/vertical-allrepos/yaylog.svg)](https://repology.org/project/yaylog/versions) ![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/Zweih/yaylog/total?style=for-the-badge&logo=archlinux&label=Downloads%20Since%202%2F4%2F2025&color=%20%231793d0)

![Alt](https://repobeats.axiom.co/api/embed/7a20b73b689d45d678001c582a9d1f124dca31ba.svg "Repobeats analytics image")


this package is compatible with the following distributions:
 - [arch linux](https://archlinux.org)
 - [manjaro](https://manjaro.org/)
 - [steamOS](https://store.steampowered.com/steamos)
 - [garuda linux](https://garudalinux.org/)
 - [endeavourOS](https://endeavouros.com/)
 - [artix linux](https://artixlinux.org/)
 - the 50 other arch-based distros, as long as it has pacman installed 

## features

- list installed packages with date/timestamps, dependencies, provisions, requirements, size on disk, and version
- filter by explicitly installed packages
- filter by packages installed as dependencies
- filter by packages required by a specific package
- filter by a specific installation date or date range
- filter by package size or size range
- filter by package name (substring match)
- sort by installation date, alphabetically, or by size on disk
- output as a table or JSON

## why is it called yaylog if it works with other AUR helpers?
because yay is my preferred AUR helper and the name has a good flow.

## is it good?
[yes.](https://news.ycombinator.com/item?id=3067434)

## roadmap

- [x] rewrite in golang
- [x] additional filters
- [ ] list possibly or confirmed stale/abandoned packages
- [x] sort by size on disk
- [ ] dependency graph
- [x] concurrent filtering
- [x] filter by size on disk
- [x] asynchronous progress bar
- [x] channel-based aggregation
- [x] concurrent sorting
- [x] search by text input
- [x] package versions
- [x] filter by date range
- [x] concurrent file reading (2x speed boost)
- [x] remove expac as a dependency (3x speed boost)
- [x] package provisions
- [x] optional full timestamp 
- [x] add CI to release binaries
- [x] remove go as a dependency
- [x] filter by range of size on disk
- [x] user defined columns
- [x] dependencies of each package
- [x] reverse-dependencies of each package (required-by field)
- [ ] package descriptions
- [ ] package URLs
- [ ] package architecture
- [ ] name exclusion filter
- [ ] self-referencing column
- [x] JSON output
- [x] no-headers option
- [ ] provides filter
- [ ] depends filter
- [x] all-columns option
- [x] required-by filter
- [ ] key/value output
- [ ] list of packages in required-by filter
- [x] config dependency injection for testing
- [ ] more extensive testing
- [x] metaflag for all filters

## installation

### from AUR (**recommended**)
install using [AUR helper](https://wiki.archlinux.org/title/AUR_helpers) like `yay`:
```bash
yay -S yaylog
```

if you prefer to install a pre-compiled binary* using the AUR, use the `yaylog-bin` package instead.

***note**: binaries are automatically, securely, and transparently compiled with github CI when a version release is created. you can audit the binary creation by checking the relevant github action for each release version.

for the latest (unstable) version from git w/ the AUR, use `yaylog-git`*.  

***note**: this is not recommended for most users 


### building from source + manual installation
**note**: this packages is specific to arch-based linux distributions

1. clone the repo:
   ```bash
   git clone https://github.com/zweih/yaylog.git
   cd yaylog
   ```
2. build the binary:
   ```bash
   go build -o yaylog ./cmd/yaylog
   ```
3. copy the binary to your system's `$PATH`:
   ```bash
   sudo install -m755 yaylog /usr/bin/yaylog
   ```
4. copy the manpage:
   ```bash
   sudo install -m644 yaylog.1 /usr/share/man/man1/yaylog.1
   ```

## usage

```bash
yaylog [options]
```

### options
- `-n <number>` | `--number <number>`: number of recent packages to display (default: 20)
- `-a` | `all`: show all installed packages (ignores `-n`)
- `-f <filter>` | `--filter <field>=<value>`: apply multiple filters for a flexible query system. Can be used multiple times. example:
  - `--filter size=10MB:1GB` -> filter by size range
  - `--filter date=2024-01-01:` -> filter by installation date
  - `--filter name=firefox` -> filter by package name
  - `--filter reason=explicit` -> filter by explicit installations
  - `--filter required-by=vlc` -> show packages required by VLC
- `--sort <mode>`: sort results by:
  - `date` (default) - sort by installation date
  - `alphabetical` - sort alphabetically by package name
  - `size:asc` / `size:desc` - sort by package size (ascending or descending)
- `--no-headers`: omit column headers in table output (useful for scripting)
- `--columns <list>`: comma-separated list of columns to display (cannot use with `--all-columns` or `--add-columns`)
- `--add-columns <list>`: comma-separated list of columns to add to defaults or `all-columns`
- `--all-columns`: show all available columns in the output (overrides defaults)
- `--full-timestamp`: display the full timestamp (date and time) of package installations instead of just the date
- `--json`: output results in JSON format (overrides table output and `--full-timestamp`)
- `--no-progress`: force no progress bar outside of non-interactive environments
- `--required-by <package-name>`: show only packages that are required by the specified package
- `-h` | `--help`: print help info

### filtering with `--filter`
the `--filter` (short-flag: `-f`) flag allow for powerful filtering of installed packages. filters can be combined by using multiple filter flags. 

short-flag filters and long-flag filters can be combined.

#### available filters
| filter type  | syntax | description |
|-------------|--------|-------------|
| **date** | `date=<value>` | filters by installation date. Supports exact dates, ranges (`YYYY-MM-DD:YYYY-MM-DD`), and open-ended ranges (`YYYY-MM-DD:` or `:YYYY-MM-DD`) |
| **required by** | `required-by=<package>` | shows only packages that are required by a specific package |
| **nme** | `name=<package>` |  filters by package name (substring match) |
| **installation reason** | `reason=explicit` / `reason=dependencies` | filters packages by installation reason: explicitly installed or installed as a dependency |
| **size** | `size=<value>` | filters by package size on disk. supports exact values (`10MB`), ranges (`10MB:1GB`), and open-ended ranges (`:500KB`, `1GB:`) |

### available columns
- `date` - installation date of the package
- `name` - package name
- `reason` - installation reason (explicit/dependency)
- `size` - package size on disk
- `version` - installed package version
- `depends` - list of dependencies (output can be long)
- `required-by` - list of packages required by the package and are dependent (output can be long) 
- `provides` - list of alternative package names or shared libraries provided by package (output can be long)

### JSON output
the `--json` flag outputs the package data as structured JSON instead of a table. this can be useful for scripts or automation.

example:
```bash
yaylog -f name=sqlite --all-columns --json
```

`sqlite` is one of the few packages that actually has all the fields populated.

output format:
```json
[
  {
    "timestamp": "2025-02-26T16:33:47Z",
    "name": "sqlite",
    "reason": "dependency",
    "size": 21074944,
    "version": "3.48.0-2",
    "depends": [
      "readline",
      "zlib",
      "glibc"
    ],
    "requiredBy": [
      "docker",
      "gnupg",
      "libsoup3",
      "nss",
      "openslide",
      "tinysparql",
      "util-linux-libs"
    ],
    "provides": [
      "sqlite3=3.48.0",
      "libsqlite3.so=0-64"
    ]
  }
]
```

### tips & tricks

- when using multiple short flags, the -n flag must be last since it consumes the next argument.
this follows standard unix-style flag parsing, where positional arguments (like numbers)
are treated as separate parameters.
  
  invalid:
  ```bash
  yaylog -ne 15  # incorrect usage 
  ```
  valid:
  ```bash
  yaylog -en 15
  ```

- the `depends`, `provides`, `required-by` columns output can be lengthy, packages like `glibc` are required by thousands of packages. to improve readability, pipe the output to `less`:
  ```bash
  yaylog --columns name,depends | less
  ```
- all options that take an argument can also be used in the `--<flag>=<argument>` format:
  ```bash
  yaylog --size=100MB:1GB --date=:2024-06-30 --number=100
  yaylog --name=gtk --sort=alphabetical
  ```
  boolean flags can also be explicitly set using `--<flag>=true` or `--<flag>=false`:
  ```bash
  yaylog --explicit=true --dependencies=false --no-progress=true
  ```
  string arguments can also be surrounded with quotes or double-quotes:
  ```bash
  yaylog --sort="alphabetical" --name="vim"
  ```

  this can be useful for scripts and automation where you might want to avoid any and all ambiguity.

  **note**: `--no-progress` is automatically set to `true` when in a non-interactive environment, so you can pipe `|` into programs like `cat`, `grep`, or `less` without issue

- the `--no-headers` flag is useful when processing output in scripts. It removes the header row, making it easier to parse package lists with tools like `awk`, `sed`, or `cut`:
  ```bash
  yaylog --no-headers --columns name,size | awk '{print $1, $2}'
  ```

### examples

 1. **show the last 10 installed packages**  
   ```bash
   yaylog -n 10
   ```
 2. **show all explicitly installed packages**  
   ```bash
   yaylog -a -f reason=explicit
   ```
 3. **show only dependencies installed on a specific date**  
   ```bash
   yaylog -f reason=dependencies -f date=2024-12-25
   ```
 4. **show all packages sorted alphabetically**  
   ```bash
   yaylog -a --sort alphabetical
   ```
 5. **show the 15 most recent explicitly installed packages**  
   ```bash
   yaylog -f reason=explicit -n 15
   ```
 6. **show packages installed between july 1, 2023, and december 31, 2023**  
   ```bash
   yaylog -f date=2023-07-01:2023-12-31
   ```
 7. **show the 20 most recently installed packages larger than 20MB**  
   ```bash
   yaylog -f size=20MB: -n 20
   ```
 8. **show all dependencies smaller than 500KB**  
   ```bash
   yaylog -f reason=dependencies -f size=:500KB
   ```
 9. **show packages between 100MB and 1GB installed up to june 30, 2024**  
   ```bash
   yaylog -f size=100MB:1GB -f date=:2024-06-30
   ```
10. **show all packages sorted by size in descending order, installed after january 1, 2024**  
   ```bash
   yaylog -a --sort size:desc -f date=2024-01-01:
   ```
11. **search for installed packages containing "python"**  
   ```bash
   yaylog -f name=python
   ```
12. **search for explicitly installed packages containing "lib" that are between 10MB and 1GB in size**  
   ```bash
   yaylog -f reason=explicit -f name=lib -f size=10MB:1GB
   ```
13. **search for packages containing "linux" installed between january 1 and june 30, 2024**  
   ```bash
   yaylog -f name=linux -f date=2024-01-01:2024-06-30
   ```
14. **search for packages containing "gtk" installed after january 1, 2023, and at least 5MB in size**  
   ```bash
   yaylog -f name=gtk -f date=2023-01-01: -f size=5MB:
   ```
15. **show packages with name, version, and size**  
   ```bash
   yaylog --columns name,version,size
   ```
16. **show package names and dependencies with `less` for readability**  
   ```bash
   yaylog --columns name,depends | less
   ```
17. **output package data in JSON format**  
   ```bash
   yaylog --json
   ```
18. **save all explicitly installed packages to a JSON file**  
   ```bash
   yaylog -f reason=explicit --json > explicit-packages.json
   ```
19. **output all packages sorted by size (descending) in JSON**  
   ```bash
   yaylog --json -a --sort size:desc
   ```
20. **output JSON with specific columns**  
   ```bash
   yaylog --json --columns name,version,size
   ```
21. **show all available package details**  
   ```bash
   yaylog --all-columns
   ```
22. **output all packages with all columns/fields in JSON format**  
   ```bash
   yaylog -a --all-columns --json
   ```
23. **show package names and sizes without headers for scripting**  
   ```bash
   yaylog --no-headers --columns name,size
   ```
24. **show all packages required by "firefox"**  
   ```bash
   yaylog -f required-by=firefox
   ```
25. **show all packages required by "gtk3" that are at least 50MB in size**  
   ```bash
   yaylog -f required-by=gtk3 -f size=50MB:
   ```
26. **show packages required by "vlc" and installed after january 1, 2024**  
   ```bash
   yaylog -f required-by=vlc -f date=2024-01-01:
   ```
