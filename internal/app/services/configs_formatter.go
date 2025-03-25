package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hiddify/ray2sing/ray2sing"
	"pingo/configs"
)

type VmessConfigsFormatter struct {
	configuration *configs.Configuration
}

func NewVmessConfigsFormatter(configuration *configs.Configuration) *VmessConfigsFormatter {
	return &VmessConfigsFormatter{
		configuration: configuration,
	}
}

func (formatter *VmessConfigsFormatter) Format(rawConfig string) (string, error) {
	outboundStr, err := ray2sing.Ray2Singbox(rawConfig, true)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(outboundStr), &data); err != nil {
		return "", err
	}

	var internalOutbounds = ""
	if outbounds, ok := data["outbounds"].([]interface{}); ok && len(outbounds) > 0 {
		if first, ok := outbounds[0].(map[string]interface{}); ok {
			if xconfig, ok := first["xconfig"].(map[string]interface{}); ok {
				if outboundsInterface, ok := xconfig["outbounds"].([]interface{}); ok {
					internalOutboundsStr, err := json.Marshal(outboundsInterface)
					if err != nil {
						return "", err
					}
					internalOutbounds = string(internalOutboundsStr)
				}
			}
		}
	}

	if internalOutbounds == "" {
		return "", errors.New("no internal outbound found")
	}

	formattedConfig := formatter.buildJson(internalOutbounds)

	return formattedConfig, nil

}

func (formatter *VmessConfigsFormatter) buildJson(internalOutbounds string) string {
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
}`, dns, inbounds, log, internalOutbounds, policy, routing, stats)
	return formattedConfig
}

type VlessConfigsFormatter struct {
	configuration *configs.Configuration
}

func NewVlessConfigsFormatter(configuration *configs.Configuration) *VlessConfigsFormatter {
	return &VlessConfigsFormatter{
		configuration: configuration,
	}
}
func (formatter *VlessConfigsFormatter) Format(rawConfig string) (string, error) {

	outboundStr, err := ray2sing.Ray2Singbox(rawConfig, true)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(outboundStr), &data); err != nil {
		return "", err
	}

	var internalOutbounds = ""
	if outbounds, ok := data["outbounds"].([]interface{}); ok && len(outbounds) > 0 {
		if first, ok := outbounds[0].(map[string]interface{}); ok {
			if xconfig, ok := first["xconfig"].(map[string]interface{}); ok {
				if outboundsInterface, ok := xconfig["outbounds"].([]interface{}); ok {
					internalOutboundsStr, err := json.Marshal(outboundsInterface)
					if err != nil {
						return "", err
					}
					internalOutbounds = string(internalOutboundsStr)
				}
			}
		}
	}

	if internalOutbounds == "" {
		return "", errors.New("no internal outbound found")
	}

	formattedConfig := formatter.buildJson(internalOutbounds)

	return formattedConfig, nil

}

func (formatter *VlessConfigsFormatter) buildJson(internalOutbounds string) string {
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
}`, dns, inbounds, log, internalOutbounds, policy, routing, stats)
	return formattedConfig
}

type TrojanConfigsFormatter struct {
	configuration *configs.Configuration
}

func NewTrojanConfigsFormatter(configuration *configs.Configuration) *TrojanConfigsFormatter {
	return &TrojanConfigsFormatter{
		configuration: configuration,
	}
}

func (formatter *TrojanConfigsFormatter) Format(rawConfig string) (string, error) {
	outboundStr, err := ray2sing.Ray2Singbox(rawConfig, true)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(outboundStr), &data); err != nil {
		return "", err
	}

	var internalOutbounds = ""
	if outbounds, ok := data["outbounds"].([]interface{}); ok && len(outbounds) > 0 {
		if first, ok := outbounds[0].(map[string]interface{}); ok {
			if xconfig, ok := first["xconfig"].(map[string]interface{}); ok {
				if outboundsInterface, ok := xconfig["outbounds"].([]interface{}); ok {
					internalOutboundsStr, err := json.Marshal(outboundsInterface)
					if err != nil {
						return "", err
					}
					internalOutbounds = string(internalOutboundsStr)
				}
			}
		}
	}

	if internalOutbounds == "" {
		return "", errors.New("no internal outbound found")
	}

	formattedConfig := formatter.buildJson(internalOutbounds)

	return formattedConfig, nil

}

func (formatter *TrojanConfigsFormatter) buildJson(internalOutbounds string) string {
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
}`, dns, inbounds, log, internalOutbounds, policy, routing, stats)
	return formattedConfig
}

type SsConfigsFormatter struct {
	configuration *configs.Configuration
}

func NewSsConfigsFormatter(configuration *configs.Configuration) *SsConfigsFormatter {
	return &SsConfigsFormatter{
		configuration: configuration,
	}
}

func (formatter *SsConfigsFormatter) Format(rawConfig string) (string, error) {
	outboundStr, err := ray2sing.Ray2Singbox(rawConfig, true)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(outboundStr), &data); err != nil {
		return "", err
	}

	var internalOutbounds = ""
	if outbounds, ok := data["outbounds"].([]interface{}); ok && len(outbounds) > 0 {
		if first, ok := outbounds[0].(map[string]interface{}); ok {
			if xconfig, ok := first["xconfig"].(map[string]interface{}); ok {
				if outboundsInterface, ok := xconfig["outbounds"].([]interface{}); ok {
					internalOutboundsStr, err := json.Marshal(outboundsInterface)
					if err != nil {
						return "", err
					}
					internalOutbounds = string(internalOutboundsStr)
				}
			}
		}
	}

	if internalOutbounds == "" {
		return "", errors.New("no internal outbound found")
	}

	formattedConfig := formatter.buildJson(internalOutbounds)

	return formattedConfig, nil

}

func (formatter *SsConfigsFormatter) buildJson(internalOutbounds string) string {
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
}`, dns, inbounds, log, internalOutbounds, policy, routing, stats)
	return formattedConfig
}
