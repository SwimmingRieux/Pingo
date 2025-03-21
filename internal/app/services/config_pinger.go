package services

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"pingo/configs"
	"pingo/internal/app/services/abstraction"
	"pingo/internal/domain/entities"
	"pingo/internal/domain/structs"
	"strconv"
	"strings"
	"sync"
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
	pingCmd := exec.Command("curl", "-o", "/dev/null", "-s", "-w", "%{time_connect}", domain.Domain.Address)
	pingCmd.Env = append(os.Environ(),
		fmt.Sprintf("https_proxy=http://127.0.0.1:%v", listenerPort),
		fmt.Sprintf("http_proxy=http://127.0.0.1:%v", listenerPort))
	result, err := pingCmd.Output()

	if err != nil {
		return
	}
	floatPing, err := strconv.ParseFloat(strings.TrimSpace(string(result)), 64)
	if err != nil {
		return
	}
	ping := int(floatPing)

	value, ok := domainScoresMap.Load(config.ConfigId)
	if ok {
		domainScoresMap.Store(config.ConfigId, (maxScore-ping)+value.(int))
	}
}
