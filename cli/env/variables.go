package env

var (
	// Global shared variables (should always have a valid value)
	ConfigPath      string
	ClusterPath     string
	ProjectHomePath string
	IsCustomConfig  bool

	// Local options (falgs)
	ClusterAction string
	ClusterName   string
	Local         bool

	// Global options (flags)
	DebugMode bool
)
