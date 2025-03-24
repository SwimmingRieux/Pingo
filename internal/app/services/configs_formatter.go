package services

import (
	"fmt"
	"github.com/hiddify/ray2sing/ray2sing"
	"pingo/configs"
)

type VmessConfigsFormatter struct {
	configuration *configs.Configuration
}

func (formatter *VmessConfigsFormatter) Format(rawConfig string) (string, error) {
	outbound, err := ray2sing.Ray2Singbox(rawConfig, true)
	if err != nil {
		return "", err
	}
	dns := formatter.configuration.V2.DNS
	inbounds := formatter.configuration.V2.Inbounds
	log := formatter.configuration.V2.Log
	policy := formatter.configuration.V2.Policy
	routing := formatter.configuration.V2.Routing
	stats := formatter.configuration.V2.Stats
	formattedConfig := fmt.Sprintf(`{
    "dns": %s,
    "inbounds": %s,
    "log": %s,
    "outbounds": %s,
    "policy": %s,
    "routing": %s,
    "stats": %s
}`, dns, inbounds, log, outbound, policy, routing, stats)

	return formattedConfig, nil

}

type VlessConfigsFormatter struct {
	configuration *configs.Configuration
}

func (formatter *VlessConfigsFormatter) Format(rawConfig string) (string, error) {

	outbound, err := ray2sing.Ray2Singbox(rawConfig, true)
	if err != nil {
		return "", err
	}
	dns := formatter.configuration.V2.DNS
	inbounds := formatter.configuration.V2.Inbounds
	log := formatter.configuration.V2.Log
	policy := formatter.configuration.V2.Policy
	routing := formatter.configuration.V2.Routing
	stats := formatter.configuration.V2.Stats
	formattedConfig := fmt.Sprintf(`{
    "dns": %s,
    "inbounds": %s,
    "log": %s,
    "outbounds": %s,
    "policy": %s,
    "routing": %s,
    "stats": %s
}`, dns, inbounds, log, outbound, policy, routing, stats)

	return formattedConfig, nil

}

type TrojanConfigsFormatter struct {
	configuration *configs.Configuration
}

func (formatter *TrojanConfigsFormatter) Format(rawConfig string) (string, error) {
	outbound, err := ray2sing.Ray2Singbox(rawConfig, true)
	if err != nil {
		return "", err
	}
	dns := formatter.configuration.V2.DNS
	inbounds := formatter.configuration.V2.Inbounds
	log := formatter.configuration.V2.Log
	policy := formatter.configuration.V2.Policy
	routing := formatter.configuration.V2.Routing
	stats := formatter.configuration.V2.Stats
	formattedConfig := fmt.Sprintf(`{
    "dns": %s,
    "inbounds": %s,
    "log": %s,
    "outbounds": %s,
    "policy": %s,
    "routing": %s,
    "stats": %s
}`, dns, inbounds, log, outbound, policy, routing, stats)

	return formattedConfig, nil
}

type SsConfigsFormatter struct {
	configuration *configs.Configuration
}

func (formatter *SsConfigsFormatter) Format(rawConfig string) (string, error) {

	outbound, err := ray2sing.Ray2Singbox(rawConfig, true)
	if err != nil {
		return "", err
	}
	dns := formatter.configuration.V2.DNS
	inbounds := formatter.configuration.V2.Inbounds
	log := formatter.configuration.V2.Log
	policy := formatter.configuration.V2.Policy
	routing := formatter.configuration.V2.Routing
	stats := formatter.configuration.V2.Stats
	formattedConfig := fmt.Sprintf(`{
    "dns": %s,
    "inbounds": %s,
    "log": %s,
    "outbounds": %s,
    "policy": %s,
    "routing": %s,
    "stats": %s
}`, dns, inbounds, log, outbound, policy, routing, stats)

	return formattedConfig, nil

}
