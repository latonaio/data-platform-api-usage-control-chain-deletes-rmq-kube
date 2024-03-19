package dpfm_api_caller

import (
	dpfm_api_input_reader "data-platform-api-usage-control-chain-deletes-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-usage-control-chain-deletes-rmq-kube/DPFM_API_Output_Formatter"
	"fmt"
	"strings"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
)

func (c *DPFMAPICaller) UsageControlChain(
	input *dpfm_api_input_reader.SDC,
	log *logger.Logger,
) *dpfm_api_output_formatter.UsageControlChain {

	where := strings.Join([]string{
		fmt.Sprintf("WHERE usageControlChain.UsageControlChain = \"%s\ ", input.UsageControlChain.UsageControlChain),
	}, "")

	rows, err := c.db.Query(
		`SELECT 
    	usageControlChain.UsageControlChain
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_usageControlChain_usageControlChain_data as usageControlChain 
		` + where + ` ;`)
	if err != nil {
		log.Error("%+v", err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToUsageControlChain(rows)
	if err != nil {
		log.Error("%+v", err)
		return nil
	}

	return data
}
