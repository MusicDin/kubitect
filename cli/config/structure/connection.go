package structure

type IP string

type ConnectionType string
type UserString string

const (
	local  ConnectionType = "local"
	remote                = "remote"
)

type Connection struct {
	IP   *IP
	SSH  *SSH
	Type *ConnectionType
	User *UserString
}
