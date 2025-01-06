# yaylog

`yaylog` is a simple CLI util, written in **Go** / **Golang**, for arch and arch-based linux distros to list recently installed packages.

despite the name, it's not limited to `yay` and works with any package manager that logs package installations to `/var/log/pacman.log`. so it can be used with `pacman`, `yay`, `paru`, `aura`, and even `yaourt` if you're somehow still using it.

it supports optional filters for explicitly installed packages or dependencies.

[![Packaging status](https://repology.org/badge/vertical-allrepos/yaylog.svg)](https://repology.org/project/yaylog/versions)

## features

- view recently installed packages with timestamps
- filter results by explicitly installed packages
- filter results by packages installed as dependencies
- show all installed packages with alignment for readability
- sort results by installation date or alphabetically
- filter results by a specific installation date

## why is it called yaylog if it works with other package managers?
because yay is my preferred aur helper and the name has a good flow.

## installation

### from AUR (**recommended**)
install using an AUR helper like `yay`:
```bash
yay -S yaylog
```

### building from source + manual installation
1. clone the repo:
   ```bash
   git clone https://github.com/zweih/yaylog.git
   cd yaylog
   ```
2. build the binary:
   ```bash
   go build -o yaylog ./cmd/yaylog
   ```
3. copy the binary to your system's `$path`:
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
- `--sort <mode>`: sort results by:
  - `date` (default): sort by installation date
  - `alphabetical`: sort alphabetically by package name
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

   **note**: the `-e` flag must be used before the `-n` flag as the n flag consumes the next argument.

## license

this project is licensed under the MIT license. see [license](LICENSE) for details.
