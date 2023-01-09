package lib

import (
	"git.torproject.org/pluggable-transports/snowflake.git/v2/common/task"
	"io"
	"log"
	"time"

	"git.torproject.org/pluggable-transports/snowflake.git/v2/common/event"
)

func NewProxyEventLogger(logPeriod time.Duration, output io.Writer) event.SnowflakeEventReceiver {
	logger := log.New(output, "", log.LstdFlags|log.LUTC)
	el := &logEventLogger{logPeriod: logPeriod, logger: logger}
	el.task = &task.Periodic{Interval: logPeriod, Execute: el.logTick}
	el.task.WaitThenStart()
	return el
}

type logEventLogger struct {
	inboundSum      int64
	outboundSum     int64
	connectionCount int
	logPeriod       time.Duration
	task            *task.Periodic
	logger          *log.Logger
}

var inChanSum int64
var outChanSum int64
var conChan int

func (p *logEventLogger) OnNewSnowflakeEvent(e event.SnowflakeEvent) {
	switch e.(type) {
	case event.EventOnProxyConnectionOver:
		e := e.(event.EventOnProxyConnectionOver)
		p.inboundSum += e.InboundTraffic
		p.outboundSum += e.OutboundTraffic
		inChanSum += e.InboundTraffic
		outChanSum += e.OutboundTraffic
		p.connectionCount += 1
		conChan += 1
	default:
		p.logger.Println(e.String())
	}
}

func (p *logEventLogger) logTick() error {
	inbound, inboundUnit := formatTraffic(p.inboundSum)
	outbound, outboundUnit := formatTraffic(p.outboundSum)
	p.logger.Printf("In the last %v, there were %v connections. Traffic Relayed ↑ %v %v, ↓ %v %v.\n",
		p.logPeriod.String(), p.connectionCount, inbound, inboundUnit, outbound, outboundUnit)
	p.outboundSum = 0
	p.inboundSum = 0
	p.connectionCount = 0
	return nil
}

func SendStatSum() (int64, int64, int) {
	inbound, _ := formatTraffic(inChanSum)
	outbound, _ := formatTraffic(outChanSum)
	inChanSum = 0
	outChanSum = 0

	return inbound, outbound, conChan
}

func (p *logEventLogger) Close() error {
	return p.task.Close()
}
