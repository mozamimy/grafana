package alerting

import (
	"sync"
	"time"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
)

type RuleReader interface {
	Fetch() []*Rule
}

type DefaultRuleReader struct {
	sync.RWMutex
	serverID       string
	serverPosition int
	clusterSize    int
	log            log.Logger
}

func NewRuleReader() *DefaultRuleReader {
	ruleReader := &DefaultRuleReader{
		log: log.New("alerting.ruleReader"),
	}

	go ruleReader.initReader()
	return ruleReader
}

func (arr *DefaultRuleReader) initReader() {
	heartbeat := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-heartbeat.C:
			arr.heartbeat()
		}
	}
}

func (arr *DefaultRuleReader) Fetch() []*Rule {
	cmd := &m.GetAllAlertsQuery{}

	if err := bus.Dispatch(cmd); err != nil {
		arr.log.Error("Could not load alerts", "error", err)
		return []*Rule{}
	}

	res := make([]*Rule, 0)
	for _, ruleDef := range cmd.Result {
		if model, err := NewRuleFromDBAlert(ruleDef); err != nil {
			arr.log.Error("Could not build alert model for rule", "ruleId", ruleDef.Id, "error", err)
		} else {
			res = append(res, model)
		}
	}

	return res
}

func (arr *DefaultRuleReader) heartbeat() {

	//Lets cheat on this until we focus on clustering
	//log.Info("Heartbeat: Sending heartbeat from " + this.serverId)
	arr.clusterSize = 1
	arr.serverPosition = 1

	/*
		cmd := &m.HeartBeatCommand{ServerId: this.serverId}
		err := bus.Dispatch(cmd)

		if err != nil {
			log.Error(1, "Failed to send heartbeat.")
		} else {
			this.clusterSize = cmd.Result.ClusterSize
			this.serverPosition = cmd.Result.UptimePosition
		}
	*/
}
