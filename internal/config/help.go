package config

import (
	"fmt"

	"github.com/spf13/pflag"
)

func PrintHelp() {
	const helpPart1 = `Usage:
  qp [command] [args] [options]

Commands:
  select <list> | s <list>      Fields to display (comma-separated)
                                - 'select all'      -> all fields
                                - 'select default'  -> default fields
                                - e.g. 'select default,version'

  where <query> | w <query>     Refine package results using one or more queries
                                - Supports: field=value, field==value
                                - Range: updated=2024-01-01:2024-01-10
                                - Existence: has:depends, no:conflicts
                                - Depth: depends=package@2, required-by=pkg@3

  order <field>:<dir> | o <..>  Sort by field in asc/desc order

  limit <num> | l <num>         Limit number of results (default: 20)
                                - Use 'limit all' to show everything
                                - Use end:<num> / mid:<num> prefix to display from
                                    different parts of the output

  format <type> | f <type>      Output format: table, json, kv (default: table)
                                - 'format table' -> tabular output with headers (default)
                                - 'format json'  -> JSON array output
                                - 'format kv'    -> key-value pairs (best for selecting
                                    many fields)

Options:
`

	const helpPart2 = `
Query Types:
  - String match:
      field=value       -> fuzzy match (e.g., name=gtk)
      field==value      -> exact match (e.g., name==bash)

  - Range match (dates, sizes):
      field=start:end   -> fuzzy match
      field==start:end  -> exact match
      Examples:
        size=10MB:1GB
        updated==2024-01-01
        updated=2024-01-01: (open-ended range)

  - Existence check:
      has:field         -> field must exist or be non-empty
      no:field          -> field must not exist or be empty

  - Depth querying (relation fields):
      field=value@1     -> direct relations (default)
      field=value@2     -> second-level relations
      field=value@3     -> third-level relations, etc.
        Note: optdepends and optional-for include hard dependencies
              after depth 1

  - Logical operators
      and   -> both must match
      or    -> either can match
      not   -> exclude what follows
      q, p  -> group logic: 
                'q' = opening paren '(', 'p' = closing paren ')'
                
                Remember with:
                  'q' is for *q*uery group start
                  'p' is for query group sto*p*

  Examples:
    qp w name=gtk or name=qt
    qp w not name==vim
    qp w reason=explicit and size=50MB:
    qp w depends=python@2
    qp w q name=vim or name=emacs p and not has:depends

Match Behavior:
  - Strings: fuzzy = substring match (case-insensitive)
             strict = exact match (case-insensitive)
  - Date:    fuzzy = match by day
             strict = exact timestamp
  - Size:    fuzzy = ~0.3% byte tolerance
             strict = exact byte size

Short Command Examples:
  qp s name,size w name=vim o updated:asc l 10
  qp where name=gtk
  qp w name==bash
  qp w reason=explicit and size=50MB:
  qp w q size=10MB:1GB or size==20MB p and not has:depends

Built-in Macros:
  - 'qp w orphan' is equivalent to 'qp where no:required-by and reason=dependency'
  - 'qp w superorphan' is equivalent to 'qp where no:required-by and reason=dependency and no:optional-for'
  - 'qp w heavy' is equivalent to 'qp where size=100MB:'

Tips:
  - Queries can include comma-separated values, these act a shorthand for 'or' logic:
      arch=aarch64,any
      provides=rustc,python3

  - Pipe long output to 'less' or 'moar':
      qp select name,depends | less

  - Output for scripting:
      qp --no-headers select name,size

  - JSON output:
      qp select name,version,size format json

  - Key-Value output (ideal for selecting all fields):
     qp s all f kv

  - Quote arguments with spaces or special characters:
      qp where description="for tree-sitter"

  - To group conditions, use q and p as grouping parentheses:
    qp where q name=curl or name=openssl and no:depends
      -> matches packages named curl or openssl but only if they have no dependencies

Supported Package Origins:
  - brew      -> Homebrew (macOS/Linux)
  - deb       -> APT/dpkg (Debian, Ubuntu)
  - flatpak   -> Flatpak (universal Linux packages)
  - npm       -> npm (Global Node.js packages)
  - opkg      -> opkg (OpenWrt, embedded systems)
  - pacman    -> pacman (Arch Linux, Manjaro)
  - pipx      -> pipx (Global virtual Python applications)
  - rpm       -> RPM (Fedora, RHEL)

Default Behavior:
  - 20 results shown unless limit is specified
  - Updated, Name, Reason, and Size fields are displayed

Use 'man qp' to see all available fields
    for select, where, and order.

See full docs at https://github.com/Zweih/qp for:
  - Available fields
  - Query examples
  - JSON schema
`

	fmt.Print(helpPart1)
	pflag.PrintDefaults()
	fmt.Print(helpPart2)
}
