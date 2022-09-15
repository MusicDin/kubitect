package modelconfig

type Bridge string
type CIDR string
type Gateway string
type Mode string

type Network struct {
	Bridge  *Bridge
	CIDR    *CIDR
	Gateway *Gateway
	Mode    *Mode
}
