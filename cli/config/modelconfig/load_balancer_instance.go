package modelconfig

type InstanceId uint
type Priority uint

type LoadBalancerInstance struct {
	CPU          *CpuSize
	Host         *HostName
	Id           *InstanceId
	IP           *IP
	MAC          *MAC
	MainDiskSize *MB
	Priority     *Priority
	RAM          *MB
}
