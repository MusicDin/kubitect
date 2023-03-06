package modelconfig

import (
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	"github.com/MusicDin/kubitect/pkg/utils/validation"
)

type MasterDefault struct {
	CPU          VCpu       `yaml:"cpu"`
	RAM          GB         `yaml:"ram"`
	MainDiskSize GB         `yaml:"mainDiskSize"`
	Labels       Labels     `yaml:"labels,omitempty"`
	Taints       []Taint    `yaml:"taints,omitempty"`
	DataDisks    []DataDisk `yaml:"dataDisks,omitempty"`
}

func (d MasterDefault) Validate() error {
	return validation.Struct(&d,
		validation.Field(&d.CPU),
		validation.Field(&d.RAM),
		validation.Field(&d.MainDiskSize),
		validation.Field(&d.Labels),
		validation.Field(&d.Taints),
		validation.Field(&d.DataDisks, validation.OmitEmpty(), validation.UniqueField("Name")),
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
	return validation.Struct(&m,
		validation.Field(&m.Default),
		validation.Field(&m.Instances,
			validation.MinLen(1).Error("At least one master instance must be configured."),
			validation.Fail().When(len(m.Instances)%2 == 0).Error("Number of master instances must be odd (1, 3, 5 etc.)."),
			validation.UniqueField("Id"),
			validation.Custom(LB_REQUIRED),
		),
	)
}

func (m *Master) SetDefaults() {
	for i := range m.Instances {
		m.Instances[i].CPU = defaults.Default(m.Instances[i].CPU, m.Default.CPU)
		m.Instances[i].RAM = defaults.Default(m.Instances[i].RAM, m.Default.RAM)
		m.Instances[i].MainDiskSize = defaults.Default(m.Instances[i].MainDiskSize, m.Default.MainDiskSize)
		m.Instances[i].DataDisks = append(m.Default.DataDisks, m.Instances[i].DataDisks...)
	}
}

type MasterInstance struct {
	Name         string     `yaml:"name,omitempty" opt:"-"`
	Id           string     `yaml:"id" opt:",id"`
	Host         string     `yaml:"host,omitempty"`
	IP           IPv4       `yaml:"ip,omitempty"`
	MAC          MAC        `yaml:"mac,omitempty"`
	CPU          VCpu       `yaml:"cpu"`
	RAM          GB         `yaml:"ram"`
	MainDiskSize GB         `yaml:"mainDiskSize"`
	DataDisks    []DataDisk `yaml:"dataDisks,omitempty"`
	Labels       Labels     `yaml:"labels,omitempty"`
	Taints       []Taint    `yaml:"taints,omitempty"`
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
	defer validation.RemoveCustomValidator(VALID_POOL)

	validation.RegisterCustomValidator(VALID_POOL, poolNameValidator(i.Host))

	return validation.Struct(&i,
		validation.Field(&i.Id, validation.NotEmpty(), validation.AlphaNumericHypUS()),
		validation.Field(&i.Host, validation.OmitEmpty(), validation.Custom(VALID_HOST)),
		validation.Field(&i.IP, validation.OmitEmpty(), validation.Custom(IP_IN_CIDR)),
		validation.Field(&i.MAC, validation.OmitEmpty()),
		validation.Field(&i.CPU),
		validation.Field(&i.RAM),
		validation.Field(&i.MainDiskSize),
		validation.Field(&i.DataDisks, validation.OmitEmpty(), validation.UniqueField("Name")),
		validation.Field(&i.Labels),
		validation.Field(&i.Taints),
	)
}
