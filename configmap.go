package main

import (
	"sync"
)

type ConfigurationMap map[string]string
type NamedConfigurationsMap map[string]ConfigurationMap

type ConfigurationManager struct {
	configMap             NamedConfigurationsMap
	mutex                 sync.RWMutex
	namespacesParamsNames map[string][]string
	namespaces            []string
}

func NewConfigurationManager(paths ...string) *ConfigurationManager {
	configManager := &ConfigurationManager{}
	configManager.configMap = make(NamedConfigurationsMap)
	configManager.namespacesParamsNames = make(map[string][]string)
	configManager.namespaces = make([]string, 0)
	configManager.mutex = sync.RWMutex{}

	for _, path := range paths {
		configManager.readNamespace(path)
	}
	configManager.prepareIndexes()

	return configManager
}

func (cm *ConfigurationManager) Namespaces() []string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.namespaces
}

func (cm *ConfigurationManager) ConfigKeys(namespace string) []string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.namespacesParamsNames[namespace]
}

func (cm *ConfigurationManager) ParameterValue(namespace, key string) *string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	namespaceConfigMap, ok := cm.configMap[namespace]
	if !ok {
		return nil
	}
	value, ok := namespaceConfigMap[key]
	if !ok {
		return nil
	}
	return &value
}

func (cm *ConfigurationManager) prepareIndexes() {

}

func (cm *ConfigurationManager) readNamespace(nPath string) {

}
