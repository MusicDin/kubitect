package modelconfig

type Connection struct {
	IP   *IP
	SSH  *SSH
	Type *ConnectionType
	User *UserString
}
