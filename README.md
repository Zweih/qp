# yaylog

`yaylog` is a CLI util, written in **Go** / **Golang**, for arch and arch-based linux distros to sort/filter installed packages.

despite the name, it's not limited to `yay` and works with any package manager that uses ALPM; so it can be used with `pacman`, `yay`, `paru`, `aura`, `pamac`, and even `yaourt` if you're somehow still using it.

`yaylog` supports optional filters/sorting for install date, package name, install reason (explicit/dependency), and size on disk.

[![Packaging status](https://repology.org/badge/vertical-allrepos/yaylog.svg)](https://repology.org/project/yaylog/versions)

![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/Zweih/yaylog/total?style=for-the-badge&logo=archlinux&label=Downloads%20Since%202%2F4%2F2025&color=%20%231793d0)


this package is compatible with the following distributions:
 - arch linux
 - manjaro
 - garuda linux
 - endeavourOS
 - the 50 other arch-based distros, as long as it has pacman installed 

## features

- list installed packages with timestamps
- display package versions 
- filter results by explicitly installed packages
- filter results by packages installed as dependencies
- sort results by installation date, alphabetically, or by size on disk
- filter results by a specific installation date or date range
- filter results by package size or size range
- filter results by package name (substring match)

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
- [x] list package versions
- [x] filter by date range
- [x] concurrent file reading (2x speed boost)
- [x] remove expac as a dependency (3x speed boost)
- [x] optional full timestamp 
- [x] add CI to release binaries
- [x] remove go as a dependency
- [x] filter by range of size on disk
- [x] user defined columns
- [x] list dependencies of each package
- [ ] list packages that depend on each package

## installation

### from AUR (**recommended**)
install the binary using an AUR helper like `yay`:
```bash
yay -S yaylog-bin
```

if you prefer to install from source using the AUR, use the `yaylog` package instead.

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
- `-e` | `--explicit`: show only explicitly installed packages
- `-d` | `--dependencies`: show only packages installed as dependencies
- `--date <filter>`: filter packages by installation date. Supports:
  - `YYYY-MM-DD` - show packages installed on the specified date
  - `YYYY-MM-DD:` - show packages installed on or after the date
  - `:YYYY-MM-DD` - show packages installed up to the date
  - `YYYY-MM-DD:YYYY-MM-DD` - show packages installed within a date range
- `--size <filter>`: filter packages by size on disk. Supports:
  - `10MB` - show packages exactly 10MB in size
  - `5GB:` - show packages 5GB and larger
  - `:20KB` - show packages up to 20KB
  - `1.5MB:2GB` - show packages between 1.5MB and 2GB
  - valid units: B (bytes), KB, MB, GB
- `--name <search-term>`: filter packages by name (substring match)
  - example: `gtk` matches `gtk3`, `libgtk`, etc.
- `--sort <mode>`: sort results by:
  - `date` (default) - sort by installation date
  - `alphabetical` - sort alphabetically by package name
  - `size:asc` / `size:desc` - sort by package size (ascending or descending)
- `--columns <list>`: comma-separated list of columns to display (overrides defaults)
- `--add-columns <list>`: comma-separated list of columns to add to defaults
- `--full-timestamp`: display the full timestamp (date and time) of package installations instead of just the date
- `--no-progress`: force no progress bar outside of non-interactive environments
- `-h` | `--help`: print help info

### available columns
- `date` - installation date of the package
- `name` - package name
- `reason` - installation reason (explicit/dependency)
- `size` - package size on disk
- `version` - installed package version
- `depends` - list of dependencies (output can be long)

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

- the `depends` column output can be lengthy. to improve readability, pipe the output to `less`:
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

  **note**: `--no-progress` is automatically set to true when in a non-interactive environment, so you can pipe `|` into programs like `cat` or `grep` without issue


### examples

 1. show the last 10 installed packages:
   ```bash
   yaylog -n 10
   ```
 2. show all explicitly installed packages:
   ```bash
   yaylog -ae
   ```
 3. show only dependencies installed on a specific date:
   ```bash
   yaylog -d --date 2024-12-25
   ```
 4. show all packages sorted alphabetically:
   ```bash
   yaylog -a --sort alphabetical
   ```
 5. show the 15 most recent explicitly installed packages:
   ```bash
   yaylog -en 15
   ```
 6. show packages installed between July 1, 2023, and December 31, 2023:
   ```bash
   yaylog --date 2023-07-01:2023-12-31
   ```
 7. show the 20 most recently installed packages larger than 20MB:
   ```bash
   yaylog --size 20MB:
   ```
 8. show all dependencies smaller than 500KB:
   ```bash
   yaylog -ad --size :500KB
   ```
 9. show packages between 100MB and 1GB installed up to June 30, 2024:
   ```bash
   yaylog --size 100MB:1GB --date :2024-06-30
   ```
10. show all packages sorted by size in descending order, installed after January 1, 2024:
   ```bash
   yaylog -a --sort size:desc --date 2024-01-01:
   ```
11. search for installed packages containing "python":
   ```bash
   yaylog --name python
   ```
12. search for explicitly installed packages containing "lib" that are between 10MB and 1GB in size:
   ```bash
   yaylog -e --name lib --size 10MB:1GB
   ```
13. search for packages containing "linux" that were installed between January 1 and June 30, 2024:
   ```bash
   yaylog --name linux --date 2024-01-01:2024-06-30
   ```
14. search for packages containing "gtk" that were installed after January 1, 2023, and are at least 5MB in size:
   ```bash
   yaylog --name gtk --date 2023-01-01: --size 5MB:
   ```
15. show packages with name, version, and size:
   ```bash
   yaylog --columns name,version,size
   ```
16. show package names and dependencies with `less` for readability:
   ```bash
   yaylog --columns name,depends | less
   ```
