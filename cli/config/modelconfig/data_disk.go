package modelconfig

type DiskName string
type PoolName string
type DiskSize uint

type DataDisk struct {
	Name *DiskName
	Pool *PoolName
	Size *DiskSize
}
