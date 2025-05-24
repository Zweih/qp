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
                                - 'select all'      → all fields
                                - 'select default'  → default fields
                                - e.g. 'select default,version'

  where <query> | w <query>     Refine package results using one or more queries
                                - Supports: field=value, field==value
                                - Range: date=2024-01-01:2024-01-10
                                - Existence: has:depends, no:conflicts
                                - You can use multiple where clauses

  order <field>:<dir> | o <..>  Sort by field in asc/desc order
                                - Fields: date, build-date, name, size, license, pkgbase

  limit <num> | l <num>         Limit number of results (default: 20)
                                - Use 'limit all' to show everything
                                - Use end:<num> / mid:<num> modifier to display from
                                    different parts of the output
Options:
`

	const helpPart2 = `
Query Types:
  - String match:
      field=value       -> fuzzy match (e.g., name=gtk)
      field==value      -> exact match (e.g., name==bash)

  - Range match (date, size):
      field=start:end   -> fuzzy match
      field==start:end  -> exact match
      Examples:
        size=10MB:1GB
        date==2024-01-01
        date=2024-01-01: (open-ended range)

  - Existence check:
      has:field         -> field must exist or be non-empty
      no:field          -> field must not exist or be empty

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
    qp w q name=vim or name=emacs p and not has:depends

Match Behavior:
  - Strings: fuzzy = substring match (case-insensitive)
             strict = exact match (case-insensitive)
  - Date:    fuzzy = match by day
             strict = exact timestamp
  - Size:    fuzzy = ~0.3% byte tolerance
             strict = exact byte size

Short Command Examples:
  qp s name,size w name=vim o date:asc l 10
  qp where name=gtk
  qp w name==bash
  qp w reason=explicit and size=50MB:
  qp w q size=10MB:1GB or size==20MB p and not has:depends

Build-in Macros:
  - 'qp w orphan' is equivalent to 'qp where no:required-by and reason=dependency'
  - 'qp w superorphan' is equivalent to 'qp where no:required-by and reason=dependency and no:optional-for'

Tips:
  - Queries can include comma-separated values, these act a shorthand for 'or' logic:
      arch=aarch64,any
      provides=rustc,python3

  - Pipe long output to 'less' or 'moar':
      qp select name,depends | less

  - Output for scripting:
      qp --no-headers select name,size

  - JSON output:
      qp select name,version,size --output=json

  - Key-Value output (ideal for selecting all fields):
     qp s all --output=kv

  - Quote arguments with spaces or special characters:
      qp where description="for tree-sitter"

  - To group conditions, use q and p (like brackets):
    qp where q name=curl or name=openssl and no:depends
      -> matches packages named curl or openssl but only if they have no dependencies

Default Behavior:
  - 20 results shown unless limit is specified
  - Progress bar disabled in non-interactive terminals

Use 'man qp' to see all available fields
    for select, where, and order.

See full docs for:
  - Available fields
  - Query examples
  - JSON schema
`

	fmt.Print(helpPart1)
	pflag.PrintDefaults()
	fmt.Print(helpPart2)
}
