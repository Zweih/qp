package rpm

import (
	"fmt"
	"os"
	"os/exec"
	"qp/internal/consts"
	"strings"
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

	dnfQuery := `
    SELECT
      r.name, ti.reason 
	  FROM
      trans_item ti 
	  JOIN
      item i ON ti.item_id = i.id 
	  JOIN
      rpm r ON i.id = r.item_id 
	  WHERE
      i.item_type = 0 AND ti.action = 1;
  `

	cmd := exec.Command("sqlite3", historyPath, dnfQuery)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return parseHistoryOutput(string(output))
	}

	yumQuery := `
    SELECT 
      p.name, 
	    CASE 
	      WHEN py.yumdb_val = 'user' THEN 2
	      WHEN py.yumdb_val = 'dep' THEN 1  
	      WHEN tdp.state IN ('Install', 'True-Install') THEN 2
	      WHEN tdp.state = 'Dep-Install' THEN 1 ELSE 0
	    END as reason
	  FROM
      pkgtups p
	  JOIN
      trans_data_pkgs tdp ON p.pkgtupid = tdp.pkgtupid
	  LEFT JOIN
      pkg_yumdb py ON p.pkgtupid = py.pkgtupid AND py.yumdb_key = 'reason'
	  WHERE
      tdp.state IN ('Install', 'True-Install', 'Dep-Install')
	  GROUP BY
      p.name;
  `

	cmd = exec.Command("sqlite3", historyPath, yumQuery)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return map[string]string{}, fmt.Errorf("both DNF and YUM queries failed: %v", err)
	}

	return parseHistoryOutput(string(output))
}

func parseHistoryOutput(output string) (map[string]string, error) {
	reasonMap := make(map[string]string)

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		reasonCode := strings.TrimSpace(parts[1])

		var pkgReason string
		switch reasonCode {
		case "1":
			pkgReason = consts.ReasonDependency
		case "2":
			pkgReason = consts.ReasonExplicit
		case "4":
			pkgReason = consts.ReasonDependency
		}

		if pkgReason != "" {
			reasonMap[name] = pkgReason
		}
	}

	return reasonMap, nil
}
