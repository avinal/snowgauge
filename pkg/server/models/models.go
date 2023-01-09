package models

import (
	sf "github.com/avinal/snowgauge/pkg/lib"
	"time"
)

type model struct {
}

var totalUp, totalDown int64

func (m model) GetLiveNetworkStat() (networks <-chan Network, err error) {
	ch := make(chan Network)
	go GetValues(ch)
	return ch, nil
}

func GetValues(ch chan<- Network) error {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		tin, tout, con := sf.SendStatSum()
		in, out := sf.SendCurrent()
		totalUp += tout
		totalDown += tin
		newVal := Network{UploadBytes: out, DownloadBytes: in, TotalUp: totalUp, TotalDown: totalDown, Connections: con}
		ch <- newVal
	}
}

type Model interface {
	GetLiveNetworkStat() (<-chan Network, error)
}

type Network struct {
	UploadBytes   int64 `json:"upbytes"`
	DownloadBytes int64 `json:"downbytes"`
	Connections   int   `json:"connections"`
	TotalUp       int64 `json:"totalUp"`
	TotalDown     int64 `json:"totalDown"`
}

func NewModel() Model {
	return &model{}
}
