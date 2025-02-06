package services

import "net"

type VmessTypeSetter struct {
}

func (setter *VmessTypeSetter) SetPort(listener net.Listener, configPath string) error {
	return nil

}

type VlessTypeSetter struct {
}

func (setter *VlessTypeSetter) SetPort(listener net.Listener, configPath string) error {
	return nil
}

type TrojanTypeSetter struct {
}

func (setter *TrojanTypeSetter) SetPort(listener net.Listener, configPath string) error {
	return nil
}

type SsTypeSetter struct {
}

func (setter *SsTypeSetter) SetPort(listener net.Listener, configPath string) error {
	return nil
}
