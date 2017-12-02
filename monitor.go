package main

import (
	mm "github.com/gonfigserver/mapmonitor"
	"path/filepath"
)

func NewMonitor() (mm.Monitor, error) {

	var config = mm.MonitorConfiguration{}
	config.PropertyFileMask = mm.PropertyFileMask
	config.FailOnDuplicates = *pFailOnDup
	config.DefaultNamespace = mm.DefaultNamespace

	if absPropRoot, err := filepath.Abs(*propertyRoot); err != nil {
		return nil, err
	} else {
		config.PropertiesRoot = absPropRoot
	}

	configManager := &mm.SyncMapPropertiesMonitor{}

	if err := configManager.Init(config); err != nil {
		return nil, err
	}

	return configManager, nil
}
