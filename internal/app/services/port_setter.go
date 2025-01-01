package services

import "net"

type VmessTypeSetter struct {
}

func (setter *VmessTypeSetter) SetPort(listener net.Listener, configPath string) error {

}

type VlessTypeSetter struct {
}

func (setter *VlessTypeSetter) SetPort(listener net.Listener, configPath string) error {

}

type TrojanTypeSetter struct {
}

func (setter *TrojanTypeSetter) SetPort(listener net.Listener, configPath string) error {

}

type SsTypeSetter struct {
}

func (setter *SsTypeSetter) SetPort(listener net.Listener, configPath string) error {

}
