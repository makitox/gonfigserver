package mapmonitor

type (
	KeyValuePair struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	Property struct {
		KeyValuePair
		Namespace string `json:"namespace"`
	}

	MonitorConfiguration struct {
		FailOnDuplicates bool
		PropertyFileMask string
		DefaultNamespace string
		PropertiesRoot   string
	}

	Monitor interface {
		NamespacesList() []string
		KeysList(namespace string) []string
		Get(namespace, key string) *Property
		Init(config MonitorConfiguration) error
	}
)
