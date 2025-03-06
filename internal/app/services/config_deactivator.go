package services

import (
	"fmt"
	"os/exec"
	"pingo/configs"
)

type ConfigDeactivator struct {
	configuration configs.Configuration
}

func (d *ConfigDeactivator) Deactivate(killFunc func() error) error {
	exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "disabled").Run()
	err := killFunc()
	if err != nil {
		errText := d.configuration.Errors.V2rayDeactivateError
		return fmt.Errorf("%v %w", errText, err)
	}
	return nil
}
