package modelconfig

type Instance struct {
	CPU          *CpuSize
	Host         *HostName
	ID           *InstanceId
	IP           *IP
	MAC          *MAC
	Labels       map[LabelKey]Label
	MainDiskSize *MB
	RAM          *MB
	Taints       []Taint
	DataDisks    []DataDisk
}
