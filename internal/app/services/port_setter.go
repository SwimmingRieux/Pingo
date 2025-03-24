package services

import (
	"encoding/json"
	"net"
	"os"
)

type VmessTypeSetter struct {
}

func (setter *VmessTypeSetter) SetPort(listener net.Listener, configPath string) error {
	listenerPort := listener.Addr().(*net.TCPAddr).Port
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var configMap map[string]interface{}

	err = json.Unmarshal(raw, &configMap)
	if err != nil {
		return err
	}

	if inbounds, ok := configMap["inbounds"].([]interface{}); ok {
		for _, inbound := range inbounds {
			if entry, ok := inbound.(map[string]interface{}); ok {
				if _, ok = entry["port"]; ok {
					entry["port"] = listenerPort
				}
			}
		}
	}

	modifiedData, err := json.MarshalIndent(configMap, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, modifiedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

type VlessTypeSetter struct {
}

func (setter *VlessTypeSetter) SetPort(listener net.Listener, configPath string) error {

	listenerPort := listener.Addr().(*net.TCPAddr).Port
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var configMap map[string]interface{}

	err = json.Unmarshal(raw, &configMap)
	if err != nil {
		return err
	}

	if inbounds, ok := configMap["inbounds"].([]interface{}); ok {
		for _, inbound := range inbounds {
			if entry, ok := inbound.(map[string]interface{}); ok {
				if _, ok = entry["port"]; ok {
					entry["port"] = listenerPort
				}
			}
		}
	}

	modifiedData, err := json.MarshalIndent(configMap, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, modifiedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

type TrojanTypeSetter struct {
}

func (setter *TrojanTypeSetter) SetPort(listener net.Listener, configPath string) error {

	listenerPort := listener.Addr().(*net.TCPAddr).Port
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var configMap map[string]interface{}

	err = json.Unmarshal(raw, &configMap)
	if err != nil {
		return err
	}

	if inbounds, ok := configMap["inbounds"].([]interface{}); ok {
		for _, inbound := range inbounds {
			if entry, ok := inbound.(map[string]interface{}); ok {
				if _, ok = entry["port"]; ok {
					entry["port"] = listenerPort
				}
			}
		}
	}

	modifiedData, err := json.MarshalIndent(configMap, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, modifiedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

type SsTypeSetter struct {
}

func (setter *SsTypeSetter) SetPort(listener net.Listener, configPath string) error {

	listenerPort := listener.Addr().(*net.TCPAddr).Port
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var configMap map[string]interface{}

	err = json.Unmarshal(raw, &configMap)
	if err != nil {
		return err
	}

	if inbounds, ok := configMap["inbounds"].([]interface{}); ok {
		for _, inbound := range inbounds {
			if entry, ok := inbound.(map[string]interface{}); ok {
				if _, ok = entry["port"]; ok {
					entry["port"] = listenerPort
				}
			}
		}
	}

	modifiedData, err := json.MarshalIndent(configMap, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, modifiedData, 0644)
	if err != nil {
		return err
	}

	return nil
}
