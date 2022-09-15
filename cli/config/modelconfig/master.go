package modelconfig

type WorkerDefault struct {
	Labels       map[LabelKey]Label
	MainDiskSize *MB
	RAM          *MB
	Taints       []Taint
}

type Master struct {
	Default   *WorkerDefault
	Instances []Instance
}
