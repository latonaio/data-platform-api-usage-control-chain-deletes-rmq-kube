package dpfm_api_output_formatter

import (
	"database/sql"
	"fmt"
)

func ConvertToUsageControlChain(rows *sql.Rows) (*UsageControlChain, error) {
	defer rows.Close()
	usageControlChain := UsageControlChain{}
	i := 0

	for rows.Next() {
		i++
		err := rows.Scan(
			&usageControlChain.usageControlChain,
			&usageControlChain.IsMarkedForDeletion,
		)
		if err != nil {
			fmt.Printf("err = %+v \n", err)
			return &usageControlChain, err
		}

	}
	if i == 0 {
		fmt.Printf("DBに対象のレコードが存在しません。")
		return &usageControlChain, nil
	}

	return &usageControlChain, nil
}
