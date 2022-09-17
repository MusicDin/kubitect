package modelconfig

type Nodes struct {
	LoadBalancer *LoadBalancer
	Master       *Worker
	Worker       *Worker
}
