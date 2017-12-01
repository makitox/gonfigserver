package mapmonitor

import (
	"github.com/magiconair/properties"
	cl "github.com/op/go-logging"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"sync"
)

var logger = cl.MustGetLogger("mapmonitor")

type SyncMapPropertiesMonitor struct {
	configMap             Namespaces
	mutex                 sync.RWMutex
	namespacesParamsNames map[string][]string
	namespaces            []string
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

func (pm *SyncMapPropertiesMonitor) Init(path string, config MonitorConfiguration) error {
	if pm == nil {
		return errors.New("Can't initialize nil object")
	}
	pm.configMap = make(Namespaces)
	pm.namespacesParamsNames = make(map[string][]string)
	pm.namespaces = make([]string, 0)
	pm.mutex = sync.RWMutex{}
	pm.readConfigRoot(path, config)
	pm.prepareIndexes()
	return nil
}

func (pm *SyncMapPropertiesMonitor) readConfigRoot(path string, config MonitorConfiguration) error {
	if pm == nil {
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
	if err := pm.readProperties(absPropertiesRoot, DefaultNamespace, config); err != nil {
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
			if err := pm.readProperties(absNamespaceRoot, namespace, config); err != nil {
				logger.Fatal(err)
				return err
			}
		}
	}
	return nil
}

func (pm *SyncMapPropertiesMonitor) readProperties(path, namespace string, config MonitorConfiguration) error {
	logger.Infof("Read properties root:  namespace = %s, path = %s", namespace, path)
	propertyFiles, err := filepath.Glob(filepath.Join(path, config.PropertyFileMask))
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
			if pm.configMap[namespace] == nil {
				logger.Infof("\t\tNamespace %s is empty, creating one", namespace)
				pm.configMap[namespace] = Properties{}
			}
			if existValue, ok := pm.configMap[namespace][key]; ok && config.FailOnDuplicates {
				logger.Infof("\t\tproperty:  " + key + " already exist with value " + existValue)
				return errors.New("Found duplicated parameter in namespace " + namespace + " with key = " + key)
			} else {
				pm.configMap[namespace][key] = value
				//logger.Infof("property:  " + key + " = " + value)
			}

		}
	}
	return nil
}

func (pm *SyncMapPropertiesMonitor) prepareIndexes() {
	logger.Infof("Creating indexes for Configuration Manager:")
	for namespace := range pm.configMap {
		logger.Infof("\tIndex namespace: %s. Indexing keys: ", namespace)
		pm.namespaces = append(pm.namespaces, namespace)
		if pm.namespacesParamsNames[namespace] == nil {
			pm.namespacesParamsNames[namespace] = make([]string, len(pm.namespacesParamsNames[namespace]))
		}
		for paramKey := range pm.configMap[namespace] {
			logger.Infof("\t\tindex key: " + paramKey)
			pm.namespacesParamsNames[namespace] = append(pm.namespacesParamsNames[namespace], paramKey)
		}
	}
	logger.Infof("Index keylists: ")
	for _, namespace := range pm.namespaces {
		logger.Infof("\tNamespace: %s. Indexing keys: ", namespace)
		for _, k := range pm.namespacesParamsNames[namespace] {
			logger.Infof("\t\t%s", k)
		}
	}
	logger.Infof("Indexing finished")
	logger.Infof("\n")
}
