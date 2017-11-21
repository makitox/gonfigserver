package main

import (
	"github.com/magiconair/properties"
	"io/ioutil"
	"path/filepath"
	"strconv"
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
	logger.Println("GET namespaces(), going to return  ")
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.namespaces
}

func (cm *ConfigurationManager) ConfigKeys(namespace string) []string {
	logger.Println("GET ConfigKeys() with arg  " + namespace)
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.namespacesParamsNames[namespace]
}

func (cm *ConfigurationManager) ParameterValue(namespace, key string) *string {
	logger.Println("GET ParameterValue() with arg  " + namespace + " " + key)
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

func (cm *ConfigurationManager) readNamespace(nPath string) {
	absFilePath, _ := filepath.Abs(nPath)
	logger.Println("Abs path: " + absFilePath)
	files, err := ioutil.ReadDir(nPath)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Find dirs: ")
	for _, f := range files {
		if f.IsDir() {
			namespace := f.Name()
			absDirPath := filepath.Join(absFilePath, f.Name())
			logger.Print(absDirPath)

			pFiles, err := filepath.Glob(filepath.Join(absDirPath, "*.properties"))
			if err != nil {
				logger.Fatal(err)
			}
			logger.Print("Inside found " + strconv.Itoa(len(pFiles)) + " *.properties files")
			for _, pFile := range pFiles {
				logger.Println("Read: " + pFile)
				p := properties.MustLoadFile(pFile, properties.UTF8)
				for _, key := range p.Keys() {
					value, _ := p.Get(key)
					if cm.configMap[namespace] == nil {
						cm.configMap[namespace] = ConfigurationMap{}
					}
					(cm.configMap[namespace])[key] = value
					logger.Println("property:  " + key + " = " + value)
				}
			}
		}
	}
}

func (cm *ConfigurationManager) prepareIndexes() {
	logger.Println("Creating indexes")
	for namespace := range cm.configMap {
		logger.Println("Found namespace: " + namespace)
		cm.namespaces = append(cm.namespaces, namespace)
		if cm.namespacesParamsNames[namespace] == nil {
			cm.namespacesParamsNames[namespace] = make([]string, len(cm.namespacesParamsNames[namespace]))
		}
		for paramKey := range cm.configMap[namespace] {
			logger.Println(" - found param key: " + paramKey)
			cm.namespacesParamsNames[namespace] = append(cm.namespacesParamsNames[namespace], paramKey)
		}
	}
	logger.Println("Stored namespaces: ")
	for _, namespace := range cm.namespaces {
		logger.Println(namespace)
		logger.Println(" - Stored keys: ")
		for _, k := range cm.namespacesParamsNames[namespace] {
			logger.Println("    - " + k)
		}
	}
	logger.Println("Stored namespaces: ")

}
