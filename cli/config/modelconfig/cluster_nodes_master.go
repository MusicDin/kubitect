package modelconfig

import (
	"cli/utils/defaults"
	v "cli/utils/validation"
)

type MasterDefault struct {
	CPU          VCpu    `yaml:"cpu"`
	RAM          GB      `yaml:"ram"`
	MainDiskSize GB      `yaml:"mainDiskSize"`
	Labels       Labels  `yaml:"labels"`
	Taints       []Taint `yaml:"taints"`
}

func (d MasterDefault) Validate() error {
	return v.Struct(&d,
		v.Field(&d.CPU),
		v.Field(&d.RAM),
		v.Field(&d.MainDiskSize),
		v.Field(&d.Labels),
		v.Field(&d.Taints),
	)
}

func (def *MasterDefault) SetDefaults() {
	def.CPU = defaults.Default(def.CPU, defaultVCpu)
	def.RAM = defaults.Default(def.RAM, defaultRAM)
	def.MainDiskSize = defaults.Default(def.MainDiskSize, defaultMainDiskSize)
}

type Master struct {
	Default   MasterDefault    `yaml:"default"`
	Instances []MasterInstance `yaml:"instances"`
}

func (m Master) Validate() error {
	return v.Struct(&m,
		v.Field(&m.Default),
		v.Field(&m.Instances,
			v.MinLen(1).Error("At least one master instance must be configured."),
			v.Fail().When(len(m.Instances)%2 == 0).Error("Number of master instances must be odd (1, 3, 5 etc.)."),
			v.UniqueField("Id"),
			v.Custom(LB_REQUIRED),
		),
	)
}

func (m *Master) SetDefaults() {
	for i := range m.Instances {
		m.Instances[i].CPU = defaults.Default(m.Instances[i].CPU, m.Default.CPU)
		m.Instances[i].RAM = defaults.Default(m.Instances[i].RAM, m.Default.RAM)
		m.Instances[i].MainDiskSize = defaults.Default(m.Instances[i].MainDiskSize, m.Default.MainDiskSize)
	}
}

type MasterInstance struct {
	Id           string     `yaml:"id" opt:",id"`
	Host         string     `yaml:"host"`
	IP           IPv4       `yaml:"ip"`
	MAC          MAC        `yaml:"mac"`
	CPU          VCpu       `yaml:"cpu"`
	RAM          GB         `yaml:"ram"`
	MainDiskSize GB         `yaml:"mainDiskSize"`
	DataDisks    []DataDisk `yaml:"dataDisks"`
	Labels       Labels     `yaml:"labels"`
	Taints       []Taint    `yaml:"taints"`
}

func (i MasterInstance) GetTypeName() string {
	return "master"
}

func (i MasterInstance) GetID() string {
	return i.Id
}

func (i MasterInstance) GetIP() IPv4 {
	return i.IP
}

func (i MasterInstance) GetMAC() MAC {
	return i.MAC
}

func (i MasterInstance) Validate() error {
	defer v.RemoveCustomValidator(VALID_POOL)

	v.RegisterCustomValidator(VALID_POOL, poolNameValidator(i.Host))

	return v.Struct(&i,
		v.Field(&i.Id, v.NotEmpty(), v.AlphaNumericHypUS()),
		v.Field(&i.Host, v.OmitEmpty(), v.Custom(VALID_HOST)),
		v.Field(&i.IP, v.OmitEmpty(), v.Custom(IP_IN_CIDR)),
		v.Field(&i.MAC, v.OmitEmpty()),
		v.Field(&i.CPU),
		v.Field(&i.RAM),
		v.Field(&i.MainDiskSize),
		v.Field(&i.DataDisks, v.UniqueField("Name")),
		v.Field(&i.Labels),
		v.Field(&i.Taints),
	)
}
