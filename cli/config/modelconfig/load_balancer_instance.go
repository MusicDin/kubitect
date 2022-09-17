package modelconfig

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
