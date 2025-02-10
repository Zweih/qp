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

- view recently installed packages with timestamps
- filter results by explicitly installed packages
- filter results by packages installed as dependencies
- sort results by installation date or alphabetically
- filter results by a specific installation date
- sort results by size on disk

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
- [ ] channel-based aggregation
- [x] concurrent sorting
- [ ] search by text input
- [ ] list package versions
- [ ] filter by date range
- [ ] concurrent file reading
- [x] remove expac as a dependency
- [x] optional full timestamp 
- [x] add CI to release binaries
- [x] remove go as a dependency

## installation

### from AUR (**recommended**)
install using an AUR helper like `yay`:
```bash
yay -S yaylog
```

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
- `-n <number>`: number of recent packages to display (default: 20)
- `-a`: show all installed packages (ignores `-n`)
- `-e`: show only explicitly installed packages
- `-d`: show only packages installed as dependencies
- `--date <YYYY-MM-DD>`: show packages installed on the specified date
- `--size <filter>`: filter packages by size on disk
   - size filter examples:
      - `">10MB"`: show packages larger than 10MB
      - `"<500KB"`: show packages smaller than 500KB
  - quotes are required for the filter
- `--sort <mode>`: sort results by:
  - `date` (default): sort by installation date
  - `alphabetical`: sort alphabetically by package name
  - `size:asc` / `size:desc`: sort by package size on disk; ascending or descending, respectively
- `--full-timestamp`: display the full timestamp (date and time) of package installations instead of just the date
- `--no-progress`: force no progress bar outside of non-interactive environments
- `-h`: print help info

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
6. show the 20 most recently installed packages larger than 20MB:
   ```bash
   yaylog --size ">20MB"
   ```
7. show the 10 largest explicitly installed packages:
   ```bash
   yaylog -en 10 --sort size:desc
   ```
8. show all dependencies smaller than 500KB:
   ```bash
   yaylog -ad --size "<500KB"
   ```

   **note**: the `-e` flag must be used before the `-n` flag as the n flag consumes the next argument.

## license

this project is licensed under the MIT license. see [license](LICENSE) for details.
