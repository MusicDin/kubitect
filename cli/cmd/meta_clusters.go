package cmd

import (
	"cli/cluster"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type MetaClusters []cluster.ClusterMeta

func (cs MetaClusters) Names() []string {
	var names []string
	for _, c := range cs {
		names = append(names, c.Name)
	}
	return names
}

func (cs MetaClusters) FindByName(name string) *cluster.ClusterMeta {
	for _, c := range cs {
		if c.Name == name {
			return &c
		}
	}

	return nil
}

func (cs MetaClusters) CountByName(name string) int {
	var i = 0

	for _, c := range cs {
		if c.Name == name {
			i++
		}
	}

	return i
}

// AllClusters returns list of clusters meta from both global (project)
// and local clusters directory (if working directory is a Kubitect project).
func AllClusters(c *GlobalContext) (MetaClusters, error) {
	cs, err := clusters(c, false)

	if err != nil {
		return nil, err
	}

	lcs, err := clusters(c, true)

	if err == nil {
		cs = append(cs, lcs...)
	}

	return cs, nil
}

// clusters returns list of clusters meta located in either global (project)
// or local clusters directory.
func clusters(c *GlobalContext, local bool) (MetaClusters, error) {
	var path string

	if local {
		path = c.LocalClustersDir()
	} else {
		path = c.ClustersDir()
	}

	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, fmt.Errorf("failed reading cluster directory: %v", err)
	}

	var cs MetaClusters

	for _, f := range files {
		if f.IsDir() {
			name := f.Name()

			cs = append(cs, cluster.ClusterMeta{
				ClusterContext: c,
				Name:           name,
				Path:           filepath.Join(path, name),
				Local:          local,
			})
		}
	}

	return cs, nil
}
