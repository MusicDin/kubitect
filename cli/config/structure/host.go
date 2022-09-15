package structure

type HostName string

type Host struct {
	Connection           *Connection
	DataResourcePools    []DataResourcePool
	Default              bool
	Name                 *HostName
	MainResourcePoolPath *ResourcePath
}
