package rpm

import (
	"fmt"
	"os"
	"os/exec"
	"qp/internal/consts"
	"strconv"
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

	dnfQuery := fmt.Sprintf(`
    SELECT
      r.name, ti.reason 
    FROM
      trans_item ti 
    JOIN
      item i ON ti.item_id = i.id 
    JOIN
      rpm r ON i.id = r.item_id 
    WHERE
      i.item_type = 0 AND ti.action = %d;
  `, dnfActionInstall)

	cmd := exec.Command("sqlite3", historyPath, dnfQuery)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return parseHistoryOutput(string(output))
	}

	yumQuery := fmt.Sprintf(`
    SELECT 
      p.name, 
      CASE
        WHEN py.yumdb_val = '%s' THEN %d
        WHEN py.yumdb_val = '%s' THEN %d
        WHEN tdp.state IN ('%s', '%s') THEN %d
        WHEN tdp.state = '%s' THEN %d
        ELSE 0
      END as reason
    FROM
      pkgtups p
    JOIN
      trans_data_pkgs tdp ON p.pkgtupid = tdp.pkgtupid
    LEFT JOIN
      pkg_yumdb py ON p.pkgtupid = py.pkgtupid AND py.yumdb_key = 'reason'
    WHERE
      tdp.state IN ('%s', '%s', '%s')
    GROUP BY
      p.name;
  `,
		yumReasonUser, dnfReasonUser,
		yumReasonDep, dnfReasonDependency,
		yumStateInstall, yumStateTrueInstall, dnfReasonUser,
		yumStateDepInstall, dnfReasonDependency,
		yumStateInstall, yumStateTrueInstall, yumStateDepInstall)

	cmd = exec.Command("sqlite3", historyPath, yumQuery)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return map[string]string{}, fmt.Errorf(permissionError, err)
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
		reasonInt, err := strconv.Atoi(reasonCode)
		if err != nil {
			continue
		}

		var pkgReason string
		switch reasonInt {
		case dnfReasonDependency, dnfReasonWeakDep:
			pkgReason = consts.ReasonDependency
		case dnfReasonUser:
			pkgReason = consts.ReasonExplicit
		}

		if pkgReason != "" {
			reasonMap[name] = pkgReason
		}
	}

	return reasonMap, nil
}
