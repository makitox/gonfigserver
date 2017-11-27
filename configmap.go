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

func New(path string) (*ConfigurationManager, error) {
	logger.Infof("Initialization Configuration manager with init path %s", path)
	configManager := createAndInitCM()

	if err := configManager.readConfigRoot(path); err != nil {
		return nil, err
	}
	logger.Infof("\n")
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

	logger.Infof("Try property root: namespace %s path %s", DefaultNamespace, absPropertiesRoot)
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
			logger.Infof("Try property root: namespace %s path %s", namespace, absNamespaceRoot)
			if err := cm.readProperties(absNamespaceRoot, namespace, PropertyFileMask); err != nil {
				logger.Fatal(err)
				return err
			}
		}
	}
	return nil
}

func (cm *ConfigurationManager) readProperties(path, namespace, fileMask string) error {
	logger.Infof("Read properties root:  namespace = %s, path = %s", namespace, path)
	propertyFiles, err := filepath.Glob(filepath.Join(path, fileMask))
	if err != nil {
		logger.Fatal(err)
		return err
	}

	for _, propertyFile := range propertyFiles {
		logger.Infof("\tRead property file: %s", propertyFile)
		p := properties.MustLoadFile(propertyFile, properties.UTF8)
		for _, key := range p.Keys() {
			value, _ := p.Get(key)
			logger.Infof("\t\tFound property: key=%s, value=%s, namespace=%s", key, value, namespace)
			if cm.configMap[namespace] == nil {
				logger.Infof("\t\tNamespace %s is empty, creating one", namespace)
				cm.configMap[namespace] = ConfigurationMap{}
			}
			if existValue, ok := cm.configMap[namespace][key]; ok && *pFailOnDup {
				logger.Infof("\t\tproperty:  " + key + " already exist with value " + existValue)
				return errors.New("Found duplicated parameter in namespace " + namespace + " with key = " + key)
			} else {
				cm.configMap[namespace][key] = value
				//logger.Infof("property:  " + key + " = " + value)
			}

		}
	}
	return nil
}

func (cm *ConfigurationManager) prepareIndexes() {
	logger.Infof("Creating indexes for Configuration Manager:")
	for namespace := range cm.configMap {
		logger.Infof("\tIndex namespace: %s. Indexing keys: ", namespace)
		cm.namespaces = append(cm.namespaces, namespace)
		if cm.namespacesParamsNames[namespace] == nil {
			cm.namespacesParamsNames[namespace] = make([]string, len(cm.namespacesParamsNames[namespace]))
		}
		for paramKey := range cm.configMap[namespace] {
			logger.Infof("\t\tindex key: " + paramKey)
			cm.namespacesParamsNames[namespace] = append(cm.namespacesParamsNames[namespace], paramKey)
		}
	}
	logger.Infof("Index keylists: ")
	for _, namespace := range cm.namespaces {
		logger.Infof("\tNamespace: %s. Indexing keys: ", namespace)
		for _, k := range cm.namespacesParamsNames[namespace] {
			logger.Infof("\t\t%s", k)
		}
	}
	logger.Infof("Indexing finished")
	logger.Infof("\n")
}
