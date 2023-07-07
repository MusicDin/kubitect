package config

import (
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	v "github.com/MusicDin/kubitect/pkg/utils/validation"
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
	return v.Struct(&d,
		v.Field(&d.CPU),
		v.Field(&d.RAM),
		v.Field(&d.MainDiskSize),
		v.Field(&d.Labels),
		v.Field(&d.Taints),
		v.Field(&d.DataDisks, v.OmitEmpty(), v.UniqueField("Name")),
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
	return v.Struct(&w,
		v.Field(&w.Default),
		v.Field(&w.Instances, v.UniqueField("Id")),
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
	IP6          IPv6       `yaml:"ip6,omitempty"`
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
		v.Field(&i.DataDisks, v.OmitEmpty(), v.UniqueField("Name")),
		v.Field(&i.Labels),
		v.Field(&i.Taints),
	)
}
