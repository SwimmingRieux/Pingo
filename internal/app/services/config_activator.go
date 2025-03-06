package services

import (
	"fmt"
	"os/exec"
	"pingo/configs"
)

type ConfigActivator struct {
	configuration *configs.Configuration
}

func (c *ConfigActivator) Activate(path string) (func() error, error) {
	port := c.configuration.V2.DefaultPort
	host := c.configuration.V2.DefaultHost
	errText := c.configuration.Errors.ProxyVariablesSetError
	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "manual").Run(); err != nil {
		return nil, fmt.Errorf("%v %w", errText, err)
	}

	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.http", "host", host).Run(); err != nil {
		return nil, fmt.Errorf("%v %w", errText, err)
	}
	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.http", "port", fmt.Sprintf("%d", port)).Run(); err != nil {
		return nil, fmt.Errorf("%v %w", errText, err)
	}

	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.https", "host", host).Run(); err != nil {
		return nil, fmt.Errorf("%v %w", errText, err)
	}
	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.https", "port", fmt.Sprintf("%d", port)).Run(); err != nil {
		return nil, fmt.Errorf("%v %w", errText, err)
	}

	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.ftp", "host", host).Run(); err != nil {
		return nil, fmt.Errorf("%v %w", errText, err)
	}
	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.ftp", "port", fmt.Sprintf("%d", port)).Run(); err != nil {
		return nil, fmt.Errorf("%v %w", errText, err)
	}

	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.socks", "host", host).Run(); err != nil {
		return nil, fmt.Errorf("%v %w", errText, err)
	}
	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.socks", "port", fmt.Sprintf("%d", port)).Run(); err != nil {
		return nil, fmt.Errorf("%v %w", errText, err)
	}

	cmd := exec.Command("v2ray", "run", "-c", path)
	err := cmd.Start()
	if err != nil {
		errText := c.configuration.Errors.V2rayActivateError
		return nil, fmt.Errorf("%v %w", errText, err)
	}

	return cmd.Process.Kill, nil
}
