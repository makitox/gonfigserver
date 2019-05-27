package mapmonitor

import (
	"github.com/magiconair/properties"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"sync"
)

type (
	Properties map[string]string

	Namespaces map[string]Properties
)

// Basic implementation of Monitor
// based on map[][] and RWMutex
type SyncMapPropertiesMonitor struct {
	configMap             Namespaces          // holds all namespaces and theirs key-value pairs
	mutex                 sync.RWMutex        //
	namespacesParamsNames map[string][]string // keeps all keys for each namespace
	namespaces            []string            // keeps all namespaces
	monitorConfig         MonitorConfiguration
}

func (pm *SyncMapPropertiesMonitor) NamespacesList() []string {
	if pm == nil {
		return nil
	}
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return pm.namespaces
}

func (pm *SyncMapPropertiesMonitor) KeysList(namespace string) []string {
	if pm == nil {
		return nil
	}
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return pm.namespacesParamsNames[namespace]

}

func (pm *SyncMapPropertiesMonitor) Get(namespace, key string) *Property {
	if pm == nil {
		return nil
	}
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	namespaceConfigMap, ok := pm.configMap[namespace]
	if !ok {
		return nil
	}
	value, ok := namespaceConfigMap[key]
	if !ok {
		return nil
	}
	return &Property{KeyValuePair: KeyValuePair{Key: key, Value: value}, Namespace: namespace}
}

func (pm *SyncMapPropertiesMonitor) Init(config MonitorConfiguration) error {
	if pm == nil {
		return errors.New("Can't initialize nil object")
	}
	pm.monitorConfig = config
	pm.configMap = make(Namespaces)
	pm.namespacesParamsNames = make(map[string][]string)
	pm.namespaces = make([]string, 0)
	pm.mutex = sync.RWMutex{}
	pm.readConfigRoot()
	pm.prepareIndexes()
	return nil
}

func (pm *SyncMapPropertiesMonitor) readConfigRoot() error {
	if pm == nil {
		return errors.New("Try to configure empty(nil) ConfigurationManager.")
	}

	// Read property root folder as default namespace
	if err := pm.readProperties(pm.monitorConfig.PropertiesRoot, pm.monitorConfig.DefaultNamespace); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(pm.monitorConfig.PropertiesRoot)
	if err != nil {
		return err
	}
	// Read all subfolders as namespaces, where subfolder name is a namespace,
	// all files inside - parameter holders
	for _, f := range files {
		if f.IsDir() {
			namespace := f.Name()
			absNamespaceRoot := filepath.Join(pm.monitorConfig.PropertiesRoot, f.Name())
			if err := pm.readProperties(absNamespaceRoot, namespace); err != nil {
				return err
			}
		}
	}
	return nil
}

func (pm *SyncMapPropertiesMonitor) readProperties(path, namespace string) error {
	propertyFiles, err := filepath.Glob(filepath.Join(path, pm.monitorConfig.PropertyFileMask))
	if err != nil {
		return err
	}

	for _, propertyFile := range propertyFiles {
		p := properties.MustLoadFile(propertyFile, properties.UTF8)
		for _, key := range p.Keys() {
			value, _ := p.Get(key)
			if pm.configMap[namespace] == nil {
				pm.configMap[namespace] = Properties{}
			}
			if _, ok := pm.configMap[namespace][key]; ok && pm.monitorConfig.FailOnDuplicates {
				return errors.New("Found duplicated parameter in namespace " + namespace + " with key = " + key)
			} else {
				pm.configMap[namespace][key] = value
			}

		}
	}
	return nil
}

func (pm *SyncMapPropertiesMonitor) prepareIndexes() {
	for namespace := range pm.configMap {
		pm.namespaces = append(pm.namespaces, namespace)
		if pm.namespacesParamsNames[namespace] == nil {
			pm.namespacesParamsNames[namespace] = make([]string, len(pm.namespacesParamsNames[namespace]))
		}
		for paramKey := range pm.configMap[namespace] {
			pm.namespacesParamsNames[namespace] = append(pm.namespacesParamsNames[namespace], paramKey)
		}
	}
}
