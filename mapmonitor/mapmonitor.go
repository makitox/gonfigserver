package mapmonitor

type (
	Properties map[string]string

	Namespaces map[string]Properties

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
	}

	Monitor interface {
		NamespacesList() []string
		KeysList(namespace string) []string
		Get(namespace, key string) *Property
		Init(path string, config MonitorConfiguration) error
	}
)
