package poc

type Variables map[string]string

type VariableOverrides struct {
	Defaults Variables `yaml:"defaults"`
	// key is the cloud alias
	Overrides map[string]*CloudVariableOverride `yaml:"overrides"`
}

type CloudVariableOverride struct {
	Defaults Variables `yaml:"defaults"`
	// key is the deploy env
	Overrides map[string]*DeployEnvVariableOverride `yaml:"overrides"`
}

type DeployEnvVariableOverride struct {
	Defaults Variables `yaml:"defaults"`
	// key is the region name
	Overrides map[string]Variables `yaml:"overrides"`
}

type DeployEnvInfo struct {
	Cloud     string
	DeployEnv string
}

type RegionInfo struct {
	Cloud     string
	DeployEnv string
	Region    string
}
