package rpm

import (
	"database/sql"
	"fmt"
	"os"
	"qp/internal/consts"

	_ "github.com/glebarez/go-sqlite"
)

const permissionError = "failed to open history database (may need root permissions): %w\n rpm install reasons will incorrectly display as 'explicit'"

func loadInstallReasons() (map[string]string, error) {
	historyPath := findHistoryDb()
	if historyPath == "" {
		return map[string]string{}, nil
	}

	return parseRpmHistory(historyPath)
}

func parseRpmHistory(historyPath string) (map[string]string, error) {
	if _, err := os.Stat(historyPath); err != nil {
		return map[string]string{}, fmt.Errorf("history database not found: %w", err)
	}

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=ro&immutable=1", historyPath))
	if err != nil {
		return map[string]string{}, fmt.Errorf(permissionError, err)
	}
	defer db.Close()

	var count int
	testQuery := `
    SELECT
      COUNT(*)
    FROM
      sqlite_master
    WHERE
      type='table'
  `

	if err := db.QueryRow(testQuery).Scan(&count); err != nil {
		return map[string]string{}, fmt.Errorf(permissionError, err)
	}

	// dnf or yum db
	checkTableQuery := `
    SELECT 
      COUNT(*)
    FROM
      sqlite_master
    WHERE
      type='table' AND name='trans_item'
  `

	var tableExists bool
	db.QueryRow(checkTableQuery).Scan(&tableExists)

	if tableExists {
		return parseDnfHistory(db)
	}

	return parseYumHistory(db)
}

func parseDnfHistory(db *sql.DB) (map[string]string, error) {
	query := `
		SELECT 
			r.name,
			ti.reason
		FROM
      trans_item ti
		JOIN
      item i ON ti.item_id = i.id
		JOIN
      rpm r ON i.id = r.item_id  
		WHERE
      i.item_type = 0  -- rpm packages
		  AND ti.action = 1    -- install actions only
		ORDER BY
      r.name
	`

	rows, err := db.Query(query)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to query DNF history: %w", err)
	}
	defer rows.Close()

	reasonMap := make(map[string]string)

	for rows.Next() {
		var name string
		var reason int

		err := rows.Scan(&name, &reason)
		if err != nil {
			continue
		}

		var pkgReason string
		switch reason {
		case dnfReasonUser:
			pkgReason = consts.ReasonExplicit
		case dnfReasonDependency, dnfReasonWeakDep:
			pkgReason = consts.ReasonDependency
		default:
			pkgReason = ""
		}

		if pkgReason != "" {
			reasonMap[name] = pkgReason
		}
	}

	return reasonMap, nil
}

func parseYumHistory(db *sql.DB) (map[string]string, error) {
	// yum stores install reasons in both state and pkg_yumdb
	query := `
		SELECT 
			p.name,
			tdp.state,
			py.yumdb_val as reason
		FROM
      pkgtups p
		JOIN
      trans_data_pkgs tdp ON p.pkgtupid = tdp.pkgtupid
		LEFT JOIN
      pkg_yumdb py ON p.pkgtupid = py.pkgtupid AND py.yumdb_key = 'reason'
		WHERE
      tdp.state IN ('Install', 'True-Install', 'Dep-Install')  -- all install ops
		GROUP BY
      p.name  -- unique packages
		ORDER
      BY p.name
	`

	rows, err := db.Query(query)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to query YUM history: %w", err)
	}
	defer rows.Close()

	reasonMap := make(map[string]string)

	for rows.Next() {
		var name, state string
		var reason sql.NullString

		err := rows.Scan(&name, &state, &reason)
		if err != nil {
			continue
		}

		var pkgReason string

		if reason.Valid {
			switch reason.String {
			case yumReasonUser:
				pkgReason = consts.ReasonExplicit
			case yumReasonDep:
				pkgReason = consts.ReasonDependency
			}
		}

		if pkgReason == "" {
			switch state {
			case yumStateInstall, yumStateTrueInstall:
				pkgReason = consts.ReasonExplicit
			case yumStateDepInstall:
				pkgReason = consts.ReasonDependency
			}
		}

		if pkgReason != "" {
			reasonMap[name] = pkgReason
		}
	}

	return reasonMap, nil
}
