package modelconfig

type Network struct {
	Bridge  *Bridge
	CIDR    *CIDR
	Gateway *Gateway
	Mode    *Mode
}
