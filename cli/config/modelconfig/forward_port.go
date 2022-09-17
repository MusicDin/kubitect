package modelconfig

type ForwardPort struct {
	Name       *ForwardPortName
	Port       *Port
	TargetPort *Port
	Target     *PortForwardTarget
}
