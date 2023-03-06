package modelconfig

import (
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	"github.com/MusicDin/kubitect/pkg/utils/validation"
)

type WorkerDefault struct {
	CPU          VCpu       `yaml:"cpu"`
	RAM          GB         `yaml:"ram"`
	MainDiskSize GB         `yaml:"mainDiskSize"`
	Labels       Labels     `yaml:"labels,omitempty"`
	Taints       []Taint    `yaml:"taints,omitempty"`
	DataDisks    []DataDisk `yaml:"dataDisks,omitempty"`
}

func (d WorkerDefault) Validate() error {
	return validation.Struct(&d,
		validation.Field(&d.CPU),
		validation.Field(&d.RAM),
		validation.Field(&d.MainDiskSize),
		validation.Field(&d.Labels),
		validation.Field(&d.Taints),
		validation.Field(&d.DataDisks, validation.OmitEmpty(), validation.UniqueField("Name")),
	)
}

func (def *WorkerDefault) SetDefaults() {
	def.CPU = defaults.Default(def.CPU, defaultVCpu)
	def.RAM = defaults.Default(def.RAM, defaultRAM)
	def.MainDiskSize = defaults.Default(def.MainDiskSize, defaultMainDiskSize)
}

type Worker struct {
	Default   WorkerDefault    `yaml:"default"`
	Instances []WorkerInstance `yaml:"instances,omitempty"`
}

func (w Worker) Validate() error {
	return validation.Struct(&w,
		validation.Field(&w.Default),
		validation.Field(&w.Instances, validation.UniqueField("Id")),
	)
}

func (w *Worker) SetDefaults() {
	for i := range w.Instances {
		w.Instances[i].CPU = defaults.Default(w.Instances[i].CPU, w.Default.CPU)
		w.Instances[i].RAM = defaults.Default(w.Instances[i].RAM, w.Default.RAM)
		w.Instances[i].MainDiskSize = defaults.Default(w.Instances[i].MainDiskSize, w.Default.MainDiskSize)
		w.Instances[i].DataDisks = append(w.Default.DataDisks, w.Instances[i].DataDisks...)
	}
}

type WorkerInstance struct {
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

func (i WorkerInstance) GetTypeName() string {
	return "worker"
}

func (i WorkerInstance) GetID() string {
	return i.Id
}

func (i WorkerInstance) GetIP() IPv4 {
	return i.IP
}

func (i WorkerInstance) GetMAC() MAC {
	return i.MAC
}

func (i WorkerInstance) Validate() error {
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
