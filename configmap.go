package main

import (
	"github.com/magiconair/properties"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"sync"
)

type ConfigurationMap map[string]string

type NamedConfigurationsMap map[string]ConfigurationMap

type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Parameter struct {
	Pair
	Namespace string `json:"namespace"`
}

type ConfigurationManager struct {
	configMap             NamedConfigurationsMap
	mutex                 sync.RWMutex
	namespacesParamsNames map[string][]string
	namespaces            []string
}

const (
	DefaultNamespace = "default"
	PropertyFileMask = "*.properties"
)

func New(paths string) (*ConfigurationManager, error) {
	configManager := createAndInitCM()

	if err := configManager.readConfigRoot(paths); err != nil {
		return nil, err
	}

	configManager.prepareIndexes()
	return configManager, nil
}

func createAndInitCM() *ConfigurationManager {
	configManager := &ConfigurationManager{}
	configManager.configMap = make(NamedConfigurationsMap)
	configManager.namespacesParamsNames = make(map[string][]string)
	configManager.namespaces = make([]string, 0)
	configManager.mutex = sync.RWMutex{}
	return configManager
}

func (cm *ConfigurationManager) Namespaces() []string {
	if cm == nil {
		return nil
	}
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.namespaces
}

func (cm *ConfigurationManager) ConfigKeys(namespace string) []string {
	if cm == nil {
		return nil
	}
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.namespacesParamsNames[namespace]
}

func (cm *ConfigurationManager) ParameterValue(namespace, key string) *Parameter {
	if cm == nil {
		return nil
	}
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
	return &Parameter{Pair: Pair{Key: key, Value: value}, Namespace: namespace}
}

func (cm *ConfigurationManager) readConfigRoot(path string) error {
	if cm == nil {
		return errors.New("Try to configure empty(nil) ConfigurationManager.")
	}
	absPropertiesRoot, _ := filepath.Abs(path)
	files, err := ioutil.ReadDir(absPropertiesRoot)
	if err != nil {
		logger.Fatal(err)
		return err
	}

	// Read property root folder as 'default namespace'
	if err := cm.readProperties(absPropertiesRoot, DefaultNamespace, PropertyFileMask); err != nil {
		logger.Fatal(err)
		return err
	}

	// Read all subfolders as namespaces, where subfolder name is a namespace,
	// all files inside - parameter holders
	for _, f := range files {
		if f.IsDir() {
			namespace := f.Name()
			absNamespaceRoot := filepath.Join(absPropertiesRoot, f.Name())
			if err := cm.readProperties(absNamespaceRoot, namespace, PropertyFileMask); err != nil {
				logger.Fatal(err)
				return err
			}
		}
	}
	return nil
}

func (cm *ConfigurationManager) readProperties(path, namespace, fileMask string) error {
	propertyFiles, err := filepath.Glob(filepath.Join(path, fileMask))
	if err != nil {
		logger.Fatal(err)
		return err
	}
	for _, propertyFile := range propertyFiles {
		p := properties.MustLoadFile(propertyFile, properties.UTF8)
		for _, key := range p.Keys() {
			value, _ := p.Get(key)
			if cm.configMap[namespace] == nil {
				cm.configMap[namespace] = ConfigurationMap{}
			}
			if _, ok := cm.configMap[namespace][key]; ok && *pFailOnDup {
				return errors.New("Found duplicated parameter in namespace " + namespace + " with key = " + key)
			} else {
				cm.configMap[namespace][key] = value
				logger.Println("property:  " + key + " = " + value)
			}

		}
	}
	return nil
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
