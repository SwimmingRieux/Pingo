package services

import "github.com/hiddify/ray2sing/ray2sing"

type VmessConfigsFormatter struct{}

func (formatter *VmessConfigsFormatter) Format(rawConfig string) (string, error) {
	return ray2sing.Ray2Singbox(rawConfig, true)
}

type VlessConfigsFormatter struct{}

func (formatter *VlessConfigsFormatter) Format(rawConfig string) (string, error) {
	return ray2sing.Ray2Singbox(rawConfig, true)
}

type TrojanConfigsFormatter struct{}

func (formatter *TrojanConfigsFormatter) Format(rawConfig string) (string, error) {
	return ray2sing.Ray2Singbox(rawConfig, true)
}

type SsConfigsFormatter struct{}

func (formatter *SsConfigsFormatter) Format(rawConfig string) (string, error) {
	return ray2sing.Ray2Singbox(rawConfig, true)
}
