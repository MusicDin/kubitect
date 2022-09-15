package modelconfig

type ClusterName string

type Cluster struct {
	Name         *ClusterName
	Network      *Network
	Nodes        *Nodes
	NodeTemplate *NodeTemplate
}
