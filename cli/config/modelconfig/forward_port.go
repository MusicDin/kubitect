package modelconfig

type ForwardPortName string
type TargetName string

type ForwardPort struct {
	Name       *ForwardPortName
	Port       *Port
	TargetPort *Port
	Target     *TargetName
}
