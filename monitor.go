package main

import mm "github.com/gonfigserver/mapmonitor"

func New(path string, config mm.MonitorConfiguration) (mm.Monitor, error) {
	logger.Infof("Initialization Configuration manager with init path %s", path)
	configManager := &mm.SyncMapPropertiesMonitor{}

	if err := configManager.Init(path, config); err != nil {
		return nil, err
	}

	return configManager, nil
}
