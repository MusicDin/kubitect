package env

var (
	// Global shared variables (should always have a valid value)
	// See cmd/root/setup
	ConfigPath      string
	ClusterPath     string
	ProjectHomePath string
	IsCustomConfig  bool

	// Local options (flags)
	ClusterAction string
	ClusterName   string
	Local         bool
	AutoApprove   bool

	// Global options (flags)
	DebugMode bool
)
