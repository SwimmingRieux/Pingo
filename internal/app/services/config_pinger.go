package services

import (
	"fmt"
	"net"
	"os/exec"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/structs"
	"sync"
	"time"
)

type ConfigPinger struct {
	portSetterFactory abstraction.PortSetterFactory
	configuration     configs.Configuration
}

func (pinger *ConfigPinger) Ping(config entities.Config, domain structs.DomainWithRank,
	listener net.Listener, domainScoresMap *sync.Map) {
	portSetter, err := pinger.portSetterFactory.Fetch(config.Type)
	if err != nil {
		return
	}
	err = portSetter.SetPort(listener, config.Path)
	if err != nil {
		return
	}
	maxScore := pinger.configuration.MaxPingWaitTime
	if err := listener.Close(); err != nil {
		return
	}

	cmd := exec.Command("v2ray", "run", "-c", config.Path)
	defer cmd.Process.Kill()
	err = cmd.Start()
	if err != nil {
		return
	}

	listenerPort := listener.Addr().(*net.TCPAddr).Port
	start := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("example.com:%d", listenerPort), time.Duration(maxScore)*time.Millisecond)
	if err != nil {
		return
	}
	defer conn.Close()
	ping := int(time.Since(start).Milliseconds())

	value, ok := domainScoresMap.Load(config.ConfigId)
	if ok {
		domainScoresMap.Store(config.ConfigId, (maxScore-ping)+value.(int))
	}
}
