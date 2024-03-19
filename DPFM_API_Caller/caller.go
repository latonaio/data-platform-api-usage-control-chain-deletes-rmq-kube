package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-usage-control-chain-deletes-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-usage-control-chain-deletes-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-usage-control-chain-deletes-rmq-kube/config"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	database "github.com/latonaio/golang-mysql-network-connector"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
	"golang.org/x/xerrors"
)

type DPFMAPICaller struct {
	ctx  context.Context
	conf *config.Conf
	rmq  *rabbitmq.RabbitmqClient
	db   *database.Mysql
}

func NewDPFMAPICaller(
	conf *config.Conf, rmq *rabbitmq.RabbitmqClient, db *database.Mysql,
) *DPFMAPICaller {
	return &DPFMAPICaller{
		ctx:  context.Background(),
		conf: conf,
		rmq:  rmq,
		db:   db,
	}
}

func (c *DPFMAPICaller) AsyncDeletes(
	accepter []string,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	log *logger.Logger,
) (interface{}, []error) {
	var response interface{}
	if input.APIType == "deletes" {
		response = c.deleteSqlProcess(input, output, accepter, log)
	}

	return response, nil
}

func (c *DPFMAPICaller) deleteSqlProcess(
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	accepter []string,
	log *logger.Logger,
) *dpfm_api_output_formatter.Message {
	var usageControlChainData *dpfm_api_output_formatter.UsageControlChain
	for _, a := range accepter {
		switch a {
		case "UsageControlChain":
			h := c.usageControlChainDelete(input, output, log)
			usageControlChainData = h
			if h == nil {
				continue
			}
		}
	}

	return &dpfm_api_output_formatter.Message{
		UsageControlChain: usageControlChainData,
	}
}

func (c *DPFMAPICaller) usageControlChainDelete(
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	log *logger.Logger,
) *dpfm_api_output_formatter.UsageControlChain {
	sessionID := input.RuntimeSessionID
	usageControlChain := c.UsageControlChainRead(input, log)
	usageControlChain.IsMarkedForDeletion = input.UsageControlChain.IsMarkedForDeletion
	res, err := c.rmq.SessionKeepRequest(nil, c.conf.RMQ.QueueToSQL()[0], map[string]interface{}{"message": usageControlChain, "function": "UsageControlChainUsageControlChain", "runtime_session_id": sessionID})
	if err != nil {
		err = xerrors.Errorf("rmq error: %w", err)
		log.Error("%+v", err)
		return nil
	}
	res.Success()
	if !checkResult(res) {
		output.SQLUpdateResult = getBoolPtr(false)
		output.SQLUpdateError = "UsageControlChain Data cannot delete"
		return nil
	}

	return usageControlChain
}

func checkResult(msg rabbitmq.RabbitmqMessage) bool {
	data := msg.Data()
	d, ok := data["result"]
	if !ok {
		return false
	}
	result, ok := d.(string)
	if !ok {
		return false
	}
	return result == "success"
}

func getBoolPtr(b bool) *bool {
	return &b
}
