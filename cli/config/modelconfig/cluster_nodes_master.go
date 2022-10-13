package modelconfig

import v "cli/validation"

type MasterDefault struct {
	CPU          *VCpu    `yaml:"cpu"`
	RAM          *GB      `yaml:"ram"`
	MainDiskSize *GB      `yaml:"mainDiskSize"`
	Labels       *Labels  `yaml:"labels"`
	Taints       *[]Taint `yaml:"taints"`
}

func (d MasterDefault) Validate() error {
	return v.Struct(&d,
		v.Field(&d.CPU, v.OmitEmpty()),
		v.Field(&d.RAM, v.OmitEmpty()),
		v.Field(&d.MainDiskSize, v.OmitEmpty()),
		v.Field(&d.Labels, v.OmitEmpty()), // TODO: Is omit empty required?
		v.Field(&d.Taints),
	)
}

type Master struct {
	Default   *MasterDefault    `yaml:"default"`
	Instances *[]MasterInstance `yaml:"instances"`
}

func (m Master) Validate() error {
	return v.Struct(&m,
		v.Field(&m.Default),
		v.Field(&m.Instances,
			v.MinLen(1).Error("At least one master instance must be configured."),
			v.Fail().When(m.Instances != nil && len(*m.Instances)%2 == 0).Error("Number of master instances must be odd (1, 3, 5 etc.)."),
		),
	)
}

type MasterInstance struct {
	Id           *string     `yaml:"id" opt:",id"`
	Host         *string     `yaml:"host"`
	IP           *IP         `yaml:"ip"`
	MAC          *MAC        `yaml:"mac"`
	CPU          *VCpu       `yaml:"cpu"`
	RAM          *GB         `yaml:"ram"`
	MainDiskSize *GB         `yaml:"mainDiskSize"`
	DataDisks    *[]DataDisk `yaml:"dataDisks"`
	Labels       *Labels     `yaml:"labels"`
	Taints       *[]Taint    `yaml:"taints"`
}

func (i MasterInstance) Validate() error {
	defer v.RemoveCustomValidator(VALID_POOL)

	v.RegisterCustomValidator(VALID_POOL, poolNameValidator(i.Host))

	return v.Struct(&i,
		v.Field(&i.Id, v.Required()),
		v.Field(&i.Host, v.OmitEmpty(), v.Custom(VALID_HOST)),
		v.Field(&i.IP, v.OmitEmpty(), v.Custom(IP_IN_CIDR)),
		v.Field(&i.MAC, v.OmitEmpty()),
		v.Field(&i.CPU, v.OmitEmpty()),
		v.Field(&i.RAM, v.OmitEmpty()),
		v.Field(&i.MainDiskSize, v.OmitEmpty()),
		v.Field(&i.DataDisks),
		v.Field(&i.Labels, v.OmitEmpty()), // TODO: Is Omit empty required?
		v.Field(&i.Taints),
	)
}
