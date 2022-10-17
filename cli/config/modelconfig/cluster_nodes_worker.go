package modelconfig

import v "cli/validation"

type WorkerDefault struct {
	CPU          *VCpu    `yaml:"cpu"`
	RAM          *GB      `yaml:"ram"`
	MainDiskSize *GB      `yaml:"mainDiskSize"`
	Labels       *Labels  `yaml:"labels"`
	Taints       *[]Taint `yaml:"taints"`
}

func (d WorkerDefault) Validate() error {
	return v.Struct(&d,
		v.Field(&d.CPU),
		v.Field(&d.RAM),
		v.Field(&d.MainDiskSize),
		v.Field(&d.Labels),
		v.Field(&d.Taints),
	)
}

type Worker struct {
	Default   *WorkerDefault    `yaml:"default"`
	Instances *[]WorkerInstance `yaml:"instances"`
}

func (w Worker) Validate() error {
	return v.Struct(&w,
		v.Field(&w.Default),
		v.Field(&w.Instances, v.OmitEmpty(), v.UniqueField("Id")),
	)
}

func (w Worker) IPs() []string {
	if w.Instances == nil {
		return nil
	}

	var ips []string

	for _, i := range *w.Instances {
		if i.IP != nil {
			ips = append(ips, string(*i.IP))
		}
	}

	return ips
}

func (w Worker) MACs() []string {
	if w.Instances == nil {
		return nil
	}

	var macs []string

	for _, i := range *w.Instances {
		if i.MAC != nil {
			macs = append(macs, string(*i.MAC))
		}
	}

	return macs
}

type WorkerInstance struct {
	Id           *string     `yaml:"id" opt:",id"`
	Host         *string     `yaml:"host"`
	IP           *IPv4       `yaml:"ip"`
	MAC          *MAC        `yaml:"mac"`
	CPU          *VCpu       `yaml:"cpu"`
	RAM          *GB         `yaml:"ram"`
	MainDiskSize *GB         `yaml:"mainDiskSize"`
	DataDisks    *[]DataDisk `yaml:"dataDisks"`
	Labels       *Labels     `yaml:"labels"`
	Taints       *[]Taint    `yaml:"taints"`
}

func (i WorkerInstance) Validate() error {
	defer v.RemoveCustomValidator(VALID_POOL)

	v.RegisterCustomValidator(VALID_POOL, poolNameValidator(i.Host))

	return v.Struct(&i,
		v.Field(&i.Id, v.Required()),
		v.Field(&i.Host, v.OmitEmpty(), v.Custom(VALID_HOST)),
		v.Field(&i.IP, v.OmitEmpty(), v.Custom(IP_IN_CIDR)),
		v.Field(&i.MAC),
		v.Field(&i.CPU),
		v.Field(&i.RAM),
		v.Field(&i.MainDiskSize),
		v.Field(&i.DataDisks, v.OmitEmpty(), v.UniqueField("Name")),
		v.Field(&i.Labels),
		v.Field(&i.Taints),
	)
}