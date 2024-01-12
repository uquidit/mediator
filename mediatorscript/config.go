package mediatorscript

type MediatorBasicConfiguration struct {
	BackendURL    string                       `json:"backend_url,omitempty" mapstructure:"backend_url"` // we need maptructure annotation so we can read yaml files
	Log           MediatorLoggingConfiguration `json:"log,omitempty"  mapstructure:"log"`
	SSLSkipVerify bool                         `json:"ssl_skip_verify,omitempty"  mapstructure:"ssl_skip_verify"`
}

type MediatorLoggingConfiguration struct {
	File  string `json:"file,omitempty"  mapstructure:"file"`
	Level string `json:"level,omitempty"  mapstructure:"level"`
}

type MediatorLegacyConfiguration struct {
	Configuration MediatorBasicConfiguration `json:"configuration,omitempty"`
	Workflows     []LegacyWorkflow           `json:"workflows,omitempty"`
}
type LegacyWorkflow struct {
	Name  string        `json:"name,omitempty"`
	Steps []LegacySteps `json:"steps,omitempty"`
}
type LegacySteps struct {
	Name   string `json:"name,omitempty"`
	Script string `json:"script,omitempty"`
}
